package services

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"slices"
	"sync"
)

var NewList []*models.Pcb
var NewSem sync.Mutex

func PMCPAlgorithm(a, b *models.Pcb) int {
	return a.Size - b.Size
}

// New debe ser bloqueado antes que Ready
// New debe ser bloqueado antes que SuspReady
// SuspReady debe ser bloqueado antes que Ready
func ToNew(pcb *models.Pcb) {
	config.Logger.Debug(fmt.Sprintf("(%d) INIT TONEW", pcb.Pid))
	UpdatePcb(pcb, "NEW")
	NewList = append(NewList, pcb)
	config.Logger.Info(fmt.Sprintf("## (%d) Se crea el proceso - Estado: NEW", pcb.Pid))

	switch config.Config.ReadyIngressAlgorithm {
	case "FIFO":
		SuspReadySem.Lock()
		utils.LockLogSuspReady(pcb.Pid, "TONEW")
		ReadySem.Lock()
		utils.LockLogReady(pcb.Pid, "TONEW")
		// si hay mas de 1 en ReadyList significa que otro proceso no se pudo mandar por falta de memoria
		if len(SuspReadyList) == 0 {
			SuspReadySem.Unlock()
			utils.UnlockLogSuspReady(pcb.Pid, "TONEW")
			NewToReady()
		} else {
			NewSem.Unlock()
			utils.UnlockLogNew(pcb.Pid, "TONEW")
			SuspReadySem.Unlock()
			utils.UnlockLogSuspReady(pcb.Pid, "TONEW")
			ReadySem.Unlock()
			utils.UnlockLogReady(pcb.Pid, "TONEW")
		}
	case "PMCP":
		slices.SortFunc(NewList, PMCPAlgorithm)

		SuspReadySem.Lock()
		utils.LockLogSuspReady(pcb.Pid, "TONEW")
		ReadySem.Lock()
		utils.LockLogReady(pcb.Pid, "TONEW")
		if len(SuspReadyList) == 0 {
			SuspReadySem.Unlock()
			utils.UnlockLogSuspReady(pcb.Pid, "TONEW")
			NewToReady()
		} else {
			NewSem.Unlock()
			utils.UnlockLogNew(pcb.Pid, "TONEW")
			SuspReadySem.Unlock()
			utils.UnlockLogSuspReady(pcb.Pid, "TONEW")
			ReadySem.Unlock()
			utils.UnlockLogReady(pcb.Pid, "TONEW")
		}
	}
}

func ExecToExit(cpu *models.CpuModule) {
	config.Logger.Debug(fmt.Sprintf("(%d) INIT EXECTOEXIT", cpu.Pcb.Pid))

	state := cpu.Pcb.CurrentState
	UpdatePcb(cpu.Pcb, "EXIT")

	config.Logger.Info(fmt.Sprintf("## (%d) - Finaliza el proceso desde estado %s", cpu.Pcb.Pid, state))

	newState := cpu.Pcb.StateMap["NEW"]
	ready := cpu.Pcb.StateMap["READY"]
	exec := cpu.Pcb.StateMap["EXEC"]
	blocked := cpu.Pcb.StateMap["BLOCKED"]
	suspBlocked := cpu.Pcb.StateMap["SUSPBLOCKED"]
	suspReady := cpu.Pcb.StateMap["SUSPREADY"]
	exit := cpu.Pcb.StateMap["EXIT"]
	config.Logger.Info(fmt.Sprintf("## (%d) - Métricas de estado: NEW (%d) (%d), READY (%d) (%d), EXEC (%d) (%d), BLOCKED (%d) (%d), SUSP. BLOCKED (%d) (%d), SUSP. READY (%d) (%d), EXIT (%d) (%d)",
		cpu.Pcb.Pid, newState.Count, newState.Time, ready.Count, ready.Time, exec.Count, exec.Time, blocked.Count, blocked.Time, suspBlocked.Count, suspBlocked.Time, suspReady.Count, suspReady.Time, exit.Count, exit.Time))

	cpu.Busy = false
}
