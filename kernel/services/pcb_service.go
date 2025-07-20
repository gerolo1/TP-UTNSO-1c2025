package services

import (
	"fmt"
	"slices"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
)

func EmptyStateMap() map[string]*models.State {
	return map[string]*models.State{
		"NEW":         {0, 0, time.Now()},
		"READY":       {0, 0, time.Now()},
		"EXEC":        {0, 0, time.Now()},
		"SUSPREADY":   {0, 0, time.Now()},
		"BLOCKED":     {0, 0, time.Now()},
		"SUSPBLOCKED": {0, 0, time.Now()},
		"EXIT":        {0, 0, time.Now()},
	}
}

// cuando se saca algo de Cpu se debe remover manualmente, ejemplo:
// cpu.Pcb = pcb
// cpu.InitDate = time.Now()
// cpu.EstNextBurst = pcb.NextBurst
func RemoveFromCurrentStateList(pcb *models.Pcb) {
	switch pcb.CurrentState {
	case "NEW":
		NewList = slices.Delete(NewList, 0, 1)
		config.Logger.Debug(fmt.Sprintf("(%d) REMOVED FROM NEWLIST", pcb.Pid))
	case "READY":
		ReadyList = slices.Delete(ReadyList, 0, 1)
		config.Logger.Debug(fmt.Sprintf("(%d) REMOVED FROM READYLIST", pcb.Pid))
	case "SUSPREADY":
		SuspReadyList = slices.Delete(SuspReadyList, 0, 1)
		config.Logger.Debug(fmt.Sprintf("(%d) REMOVED FROM SUSPREADYLIST", pcb.Pid))
	}
}

func UpdatePcb(pcb *models.Pcb, next string) {
	config.Logger.Debug(fmt.Sprintf("(%d) INIT UPDATEPCB", pcb.Pid))
	timeNow := time.Now()

	state := pcb.StateMap[next]
	state.InitDate = timeNow
	state.Count++

	if next != "NEW" {
		state = pcb.StateMap[pcb.CurrentState]

		timeDiff := timeNow.Sub(state.InitDate)
		timeDiffMilliseconds := timeDiff.Milliseconds()
		state.Time += timeDiffMilliseconds

		if config.Config.SchedulerAlgorithm == "SJF" || config.Config.SchedulerAlgorithm == "SRT" {
			if pcb.CurrentState == "EXEC" {
				pcb.PrevBurst = float64(timeDiffMilliseconds)
			}

			if next == "READY" && pcb.CurrentState != "NEW" {
				pcb.EstPrevBurst = pcb.NextBurst
				pcb.NextBurst = config.Config.Alpha*pcb.PrevBurst + (1-config.Config.Alpha)*pcb.EstPrevBurst
			}
		}
	}
	pcb.CurrentState = next
	config.Logger.Debug(fmt.Sprintf("(%d) PCB UPDATED", pcb.Pid))
}
