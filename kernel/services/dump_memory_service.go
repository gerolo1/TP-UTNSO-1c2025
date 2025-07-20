package services

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/clients"
	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
	"net/http"
	"slices"
)

func DumpBlock(pid int, pc int) {
	SuspReadySem.Lock()
	ReadySem.Lock()
	CpuSem.Lock()
	IoDevicesSem.Lock()

	device, _ := IoDevices["DUMP"]
	cpu := FindCpuByPid(pid)
	pcb := cpu.Pcb
	AddToBlocked(device, pc, 0, pcb, "DUMP")

	cpu.Busy = false

	if len(ReadyList) > 0 {
		ReadyToExec()
	}

	device.Queue = slices.Delete(device.Queue, 0, 1)
	resp := clients.DumpProcess(pid)

	IoDevicesSem.Unlock()
	if resp.StatusCode == http.StatusOK {
		if pcb.CurrentState == "BLOCKED" {
			SuspReadySem.Unlock()

			BlockedToReady(pcb)
		} else {
			SuspBlockedToSuspReady(pcb)

			SuspReadySem.Unlock()
		}
	} else {

		state := pcb.CurrentState
		UpdatePcb(pcb, "EXIT")

		config.Logger.Info(fmt.Sprintf("## (%d) - Finaliza el proceso desde estado %s", pcb.Pid, state))

		newState := pcb.StateMap["NEW"]
		ready := pcb.StateMap["READY"]
		exec := pcb.StateMap["EXEC"]
		blocked := pcb.StateMap["BLOCKED"]
		suspBlocked := pcb.StateMap["SUSPBLOCKED"]
		suspReady := pcb.StateMap["SUSPREADY"]
		exit := pcb.StateMap["EXIT"]
		config.Logger.Info(fmt.Sprintf("## (%d) - Métricas de estado: NEW (%d) (%d), READY (%d) (%d), EXEC (%d) (%d), BLOCKED (%d) (%d), SUSP. BLOCKED (%d) (%d), SUSP. READY (%d) (%d), EXIT (%d) (%d)",
			pcb.Pid, newState.Count, newState.Time, ready.Count, ready.Time, exec.Count, exec.Time, blocked.Count, blocked.Time, suspBlocked.Count, suspBlocked.Time, suspReady.Count, suspReady.Time, exit.Count, exit.Time))
		SuspReadySem.Unlock()
	}
	ReadySem.Unlock()
	CpuSem.Unlock()
}

func FindCpuByPid(pid int) *models.CpuModule {
	// Buscar en todas las listas conocidas

	for _, cpu := range CpuList {
		if cpu.Pcb.Pid == pid {
			return cpu
		}
	}

	// puedo extender con otras listas EXEC, etc.
	return nil
}
