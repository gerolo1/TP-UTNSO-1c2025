package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models/dtos"
	"github.com/sisoputnfrba/tp-golang/kernel/services"
)

func NewCpuHandler(response http.ResponseWriter, request *http.Request) {
	var body dtos.CpuHandshakeDTO

	if request.Body == nil {
		config.Logger.Error("BAD REQUEST: BODY NOT FOUND")
		http.Error(response, "BAD REQUEST: BODY NOT FOUND", http.StatusBadRequest)
		return
	}

	e := json.NewDecoder(request.Body).Decode(&body)
	if e != nil {
		config.Logger.Error("BAD REQUEST: Invalid JSON format")
		http.Error(response, "BAD REQUEST: Invalid JSON format", http.StatusBadRequest)
		return
	}
	services.NewCpu(response, body)
}

func ExitHandler(response http.ResponseWriter, request *http.Request) {
	var body dtos.ExecProcessDTO

	if request.Body == nil {
		config.Logger.Error("BAD REQUEST: BODY NOT FOUND")
		http.Error(response, "BAD REQUEST: BODY NOT FOUND", http.StatusBadRequest)
		return
	}

	e := json.NewDecoder(request.Body).Decode(&body)
	if e != nil {
		config.Logger.Error("BAD REQUEST: Invalid JSON format")
		http.Error(response, "BAD REQUEST: Invalid JSON format", http.StatusBadRequest)
		return
	}
	config.Logger.Info(fmt.Sprintf("## (%d) - Solicitó syscall: EXIT", body.Pid))

	services.ExitProcess(response, body)
}

func IoHandler(w http.ResponseWriter, r *http.Request) {
	config.Logger.Debug("INIT: CpuController - IoHandler")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		config.Logger.Debug("END: CpuController - wrong method")
		return
	}

	var dto dtos.IoRequestDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		config.Logger.Error("Invalid JSON in IoHandler")
		return
	}

	config.Logger.Info(fmt.Sprintf("## (%d) - Solicito syscall IO, Dispositivo: %s, Duracion: %d", dto.PID, dto.Device, dto.Duration))

	// Delegar la lógica al service
	services.HandleIoSyscall(dto.PID, dto.PC, dto.Device, dto.Duration)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"result":  "OK",
		"message": "Process enqueued in IO",
	})

	config.Logger.Debug("END: CpuController - IoHandler")
}
