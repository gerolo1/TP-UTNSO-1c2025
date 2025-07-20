package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models/dtos"
	"github.com/sisoputnfrba/tp-golang/kernel/services"
)

func IoCompleteHandler(w http.ResponseWriter, r *http.Request) {
	config.Logger.Debug("INIT: IoCompleteHandler")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dto dtos.IoCompleteDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	config.Logger.Debug(fmt.Sprintf("IO Complete received - Device: %s, PID: %d", dto.NameIO, dto.PID))
	services.HandleIoComplete(dto.NameIO, dto.PID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"result": "OK",
	})

	config.Logger.Debug("END: IoCompleteHandler")
}
func IoHandshakeHandler(w http.ResponseWriter, r *http.Request) {
	config.Logger.Debug("INIT: IoHandshakeHandler")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var dto dtos.IOHandshakeDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		config.Logger.Error("Invalid JSON in IoHandshakeHandler")
		return
	}

	services.RegisterIoDevice(dto.Name, dto.IP, dto.Port)

	config.Logger.Info(fmt.Sprintf("IO Module registered via handshake: %s at %s:%d", dto.Name, dto.IP, dto.Port))

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"result":  "registered",
		"message": fmt.Sprintf("IO %s registered successfully", dto.Name),
	})

	config.Logger.Debug("END: IoHandshakeHandler")
}

func IoShutdownHandler(w http.ResponseWriter, r *http.Request) {
	config.Logger.Debug("INIT: IoShutdownHandler")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var dto dtos.IOHandshakeDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		config.Logger.Error("Invalid JSON in IoHandshakeHandler")
		return
	}

	services.HandleIoShutdown(dto.Name, dto.IP, dto.Port)

	w.WriteHeader(http.StatusNoContent)
}
