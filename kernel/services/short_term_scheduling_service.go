package services

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/clients"
	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"net/http"
	"slices"
	"sync"
	"time"
)

var ReadyList []*models.Pcb
var SuspReadyList []*models.Pcb
var ReadySem sync.Mutex
var SuspReadySem sync.Mutex

func SJFAlgorithm(a, b *models.Pcb) int {
	if a.NextBurst < b.NextBurst {
		return -1
	} else if a.NextBurst > b.NextBurst {
		return 1
	}
	return 0
}

func MaxBurstCpu(a *models.CpuModule, b *models.CpuModule) int {
	timeNow := time.Now()
	timeDiffA := timeNow.Sub(a.InitDate)
	timeBurstA := a.EstNextBurst - float64(timeDiffA.Milliseconds())
	timeDiffB := timeNow.Sub(b.InitDate)
	timeBurstB := b.EstNextBurst - float64(timeDiffB.Milliseconds())

	if timeBurstA < timeBurstB {
		return -1
	} else if timeBurstA > timeBurstB {
		return 1
	}
	return 0
}

func AddToReady(pcb *models.Pcb) {
	ReadyList = append(ReadyList, pcb)

	switch config.Config.SchedulerAlgorithm {
	case "SJF":
		slices.SortFunc(ReadyList, SJFAlgorithm)
	case "SRT":
		slices.SortFunc(ReadyList, SJFAlgorithm)
	}
}

// New debe ser bloqueado antes que Ready
func NewToReady() {
	pcb := NewList[0]

	config.Logger.Debug(fmt.Sprintf("(%d) INIT NEWTOREADY", pcb.Pid))

	response := clients.AskForMemory(pcb)

	if response.StatusCode == http.StatusCreated {
		RemoveFromCurrentStateList(pcb)

		NewSem.Unlock()
		utils.UnlockLogNew(pcb.Pid, "NEWTOREADY")

		UpdatePcb(pcb, "READY")

		AddToReady(pcb)

		config.Logger.Info(fmt.Sprintf("## (%d) Pasa del estado NEW al estado READY", pcb.Pid))

		CpuSem.Lock()
		utils.LockLogCpu(pcb.Pid, "NEWTOREADY")

		ReadyToExec()

		ReadySem.Unlock()
		utils.UnlockLogReady(pcb.Pid, "NEWTOREADY")
		CpuSem.Unlock()
		utils.UnlockLogCpu(pcb.Pid, "NEWTOREADY")
	} else {
		config.Logger.Debug(fmt.Sprintf("(%d) NOT INITIALIZED: NOT ENOUGH MEMORY", pcb.Pid))

		NewSem.Unlock()
		utils.UnlockLogNew(pcb.Pid, "NEWTOREADY")
		ReadySem.Unlock()
		utils.UnlockLogReady(pcb.Pid, "NEWTOREADY")
	}
}

// New debe ser bloqueado antes que Ready
// Ready debe ser bloqueado antes de CPU
func ReadyToExec() {
	if len(ReadyList) == 0 {
		return
	}
	pcb := ReadyList[0]

	config.Logger.Debug(fmt.Sprintf("(%d) INIT READYTOEXEC", pcb.Pid))

	cpu := CheckFreeCpus()

	switch config.Config.SchedulerAlgorithm {
	case "SRT":
		if cpu != nil && !cpu.Interrupt {
			response := clients.SendToExec(pcb, cpu)

			if response.StatusCode == http.StatusOK {
				RemoveFromCurrentStateList(pcb)
				UpdatePcb(pcb, "EXEC")

				UpdateCpu(cpu, pcb)

				config.Logger.Info(fmt.Sprintf("## (%d) Pasa del estado READY al estado EXEC (CPU %s)", pcb.Pid, cpu.Id))
			}
		} else {
			cpu = slices.MaxFunc(CpuList, MaxBurstCpu)

			timeNow := time.Now()
			timeDiff := timeNow.Sub(cpu.InitDate)
			timeBurst := cpu.EstNextBurst - float64(timeDiff.Milliseconds())
			if pcb.NextBurst < timeBurst && !cpu.Interrupt {
				prevPcb := ExecToReady(pcb, cpu)

				response := clients.SendToExec(pcb, cpu)
				config.Logger.Info(fmt.Sprintf("## (%d) Pasa del estado READY al estado EXEC (CPU %s)", pcb.Pid, cpu.Id))

				if response.StatusCode == http.StatusOK {
					UpdatePcb(pcb, "EXEC")

					UpdateCpu(cpu, pcb)
					cpu.Interrupt = false

					if prevPcb.CurrentState == "READY" {
						AddToReady(prevPcb)
					}
				}
			} else {
				config.Logger.Debug(fmt.Sprintf("(%d) NOT EXECUTE", pcb.Pid))
			}
		}
	default:
		if cpu != nil {
			response := clients.SendToExec(pcb, cpu)

			if response.StatusCode == http.StatusOK {
				RemoveFromCurrentStateList(pcb)
				UpdatePcb(pcb, "EXEC")

				UpdateCpu(cpu, pcb)

				config.Logger.Info(fmt.Sprintf("## (%d) Pasa del estado READY al estado EXEC (CPU %s)", pcb.Pid, cpu.Id))
			}
		} else {
			config.Logger.Debug(fmt.Sprintf("(%d) NOT EXECUTE", pcb.Pid))
		}
	}
}

func ExecToReady(pcb *models.Pcb, cpu *models.CpuModule) *models.Pcb {
	config.Logger.Debug(fmt.Sprintf("(%d) INIT EXECTOREADY", cpu.Pcb.Pid))
	RemoveFromCurrentStateList(pcb)
	cpu.Interrupt = true

	ReadySem.Unlock()
	utils.UnlockLogReady(pcb.Pid, "EXECTOREADY")
	CpuSem.Unlock()
	utils.UnlockLogCpu(pcb.Pid, "EXECTOREADY")

	retProcess := clients.Interrupt(cpu)

	ReadySem.Lock()
	utils.LockLogReady(pcb.Pid, "EXECTOREADY")
	CpuSem.Lock()
	utils.LockLogCpu(pcb.Pid, "EXECTOREADY")

	cpu.Pcb.Pc = retProcess.PC

	if cpu.Pcb.CurrentState == "EXEC" {
		UpdatePcb(cpu.Pcb, "READY")
	}

	config.Logger.Info(fmt.Sprintf("## (%d) - Desalojado de (%s) por algoritmo SRT (Cambio de estado de EXEC a READY). Motivo: %s", cpu.Pcb.Pid, cpu.Id, retProcess.Motivo))

	return cpu.Pcb
}

func CheckSuspReady() bool {
	enoughMemory := true
	pcb := SuspReadyList[0]

	config.Logger.Debug(fmt.Sprintf("(%d) INIT CHECKSUSPREADY", pcb.Pid))

	response := clients.ResumeProcess(pcb)

	if response.StatusCode == http.StatusOK {
		RemoveFromCurrentStateList(pcb)

		UpdatePcb(pcb, "READY")

		AddToReady(pcb)
		config.Logger.Info(fmt.Sprintf("## (%d) Pasa del estado SUSP. READY al estado READY", pcb.Pid))

		ReadyToExec()
	} else {
		config.Logger.Debug(fmt.Sprintf("(%d) NOT INITIALIZED: NOT ENOUGH MEMORY TO RESUME", pcb.Pid))
		enoughMemory = false
	}
	return enoughMemory
}

func CheckNew() bool {
	enoughMemory := true
	pcb := NewList[0]

	config.Logger.Debug(fmt.Sprintf("(%d) INIT CHECKNEW", pcb.Pid))

	response := clients.AskForMemory(pcb)

	if response.StatusCode == http.StatusCreated {
		RemoveFromCurrentStateList(pcb)

		UpdatePcb(pcb, "READY")

		AddToReady(pcb)
		config.Logger.Info(fmt.Sprintf("## (%d) Pasa del estado NEW al estado READY", pcb.Pid))

		ReadyToExec()
	} else {
		config.Logger.Debug(fmt.Sprintf("(%d) NOT INITIALIZED: NOT ENOUGH MEMORY", pcb.Pid))
		enoughMemory = false
	}
	return enoughMemory
}
