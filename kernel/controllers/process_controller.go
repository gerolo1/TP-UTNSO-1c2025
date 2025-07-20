package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
	"github.com/sisoputnfrba/tp-golang/kernel/models/dtos"
	"github.com/sisoputnfrba/tp-golang/kernel/services"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
)

var LastPidSem sync.Mutex
var lastPid int = 0

func InitProcessHandler(w http.ResponseWriter, r *http.Request) {
	var body dtos.InitProcDTO

	if r.Body == nil {
		http.Error(w, "Body not found", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	LastPidSem.Lock()
	lastPid++
	pcb := &models.Pcb{
		Pid:          lastPid,
		Pc:           0,
		Size:         body.ProcSize,
		File:         body.FileName,
		StateMap:     services.EmptyStateMap(),
		CurrentState: "",
		EstPrevBurst: 0,
		PrevBurst:    0,
		NextBurst:    config.Config.InitialEstimate,
	}
	LastPidSem.Unlock()
	config.Logger.Info(fmt.Sprintf("## (%d) - Solicito syscall INIT_PROC (archivo: %s, tamaño: %d)", pcb.Pid, pcb.File, pcb.Size))

	w.WriteHeader(http.StatusNoContent)

	// Enviar al estado NEW
	go func() {
		services.NewSem.Lock()
		utils.LockLogNew(pcb.Pid, "INIT_PROC")
		services.ToNew(pcb)
	}()
	return
}
