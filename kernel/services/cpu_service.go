package services

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/clients"
	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
	"github.com/sisoputnfrba/tp-golang/kernel/models/dtos"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"

	"net/http"
	"sync"
	"time"
)

var CpuList []*models.CpuModule
var CpuSem sync.Mutex

func CheckFreeCpus() *models.CpuModule {
	for _, cpu := range CpuList {
		if !cpu.Busy && !cpu.Interrupt {
			return cpu
		}
	}
	return nil
}

func UpdateCpu(cpu *models.CpuModule, pcb *models.Pcb) {
	cpu.Pcb = pcb
	cpu.Busy = true
	cpu.InitDate = time.Now()
	cpu.EstNextBurst = pcb.NextBurst
}

// Ready debe ser bloqueado antes de CPU
func NewCpu(response http.ResponseWriter, data dtos.CpuHandshakeDTO) {
	config.Logger.Debug(fmt.Sprintf("(%s) INIT NEWCPU", data.Id))
	cpu := &models.CpuModule{data.Id, fmt.Sprintf("%s:%d", data.Ip, data.Port), 0, time.Now(), nil, false, false}

	ReadySem.Lock()
	config.Logger.Debug(fmt.Sprintf("NEWCPU (%s): READYSEM LOCKED", data.Id))
	CpuSem.Lock()
	config.Logger.Debug(fmt.Sprintf("NEWCPU (%s): CPUSEM LOCKED", data.Id))

	CpuList = append(CpuList, cpu)

	ReadySem.Unlock()
	config.Logger.Debug(fmt.Sprintf("NEWCPU (%s): READYSEM UNLOCKED", data.Id))
	CpuSem.Unlock()
	config.Logger.Debug(fmt.Sprintf("NEWCPU (%s): CPUSEM UNLOCKED", data.Id))

	response.WriteHeader(http.StatusOK)

	config.Logger.Debug(fmt.Sprintf("(%s) CPU MODULE CONNECTED TO KERNEL", data.Id))
}

// Monopoliza los recursos
func ExitProcess(response http.ResponseWriter, data dtos.ExecProcessDTO) {
	config.Logger.Debug(fmt.Sprintf("(%d) INIT EXITPROCESS", data.Pid))

	NewSem.Lock()
	utils.LockLogNew(data.Pid, "EXITPROCESS")
	SuspReadySem.Lock()
	utils.LockLogSuspReady(data.Pid, "EXITPROCESS")
	ReadySem.Lock()
	utils.LockLogReady(data.Pid, "EXITPROCESS")
	CpuSem.Lock()
	utils.LockLogCpu(data.Pid, "EXITPROCESS")

	resp := clients.ExitMemory(data.Pid)

	if resp.StatusCode == http.StatusOK {
		cpu := func() *models.CpuModule {
			for _, cpu := range CpuList {
				if cpu.Pcb.Pid == data.Pid {
					return cpu
				}
			}
			return nil
		}()

		if cpu == nil {
			config.Logger.Info(fmt.Sprintf("CPU ME DA NULL EN EXITPROCESS, PID: %d", data.Pid))
		}

		ExecToExit(cpu)

		if len(SuspReadyList) == 0 && len(NewList) == 0 && len(ReadyList) > 0 {
			ReadyToExec()
		} else {
			enoughMemory := true
			for len(SuspReadyList) > 0 && enoughMemory {
				enoughMemory = CheckSuspReady()
			}
			for len(NewList) > 0 && enoughMemory {
				enoughMemory = CheckNew()
			}
			ReadyToExec()
		}

		response.WriteHeader(http.StatusNoContent)
	}

	NewSem.Unlock()
	utils.UnlockLogNew(data.Pid, "EXITPROCESS")
	SuspReadySem.Unlock()
	utils.UnlockLogSuspReady(data.Pid, "EXITPROCESS")
	ReadySem.Unlock()
	utils.UnlockLogReady(data.Pid, "EXITPROCESS")
	CpuSem.Unlock()
	utils.UnlockLogCpu(data.Pid, "EXITPROCESS")

	config.Logger.Debug(fmt.Sprintf("(%d) ENDED EXITPROCESS", data.Pid))
}
