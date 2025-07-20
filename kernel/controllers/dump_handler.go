package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/models/dtos"
	"github.com/sisoputnfrba/tp-golang/kernel/services"
)

func DumpHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Parsear body con PID
	var recDTO dtos.DumpReceiveDTO
	err := json.NewDecoder(r.Body).Decode(&recDTO)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.Logger.Info(fmt.Sprintf("## (%d) - Solicitó syscall: DUMP", recDTO.PID))

	// 2. Bloquear proceso, hacer dump y cambiar estado
	services.DumpBlock(recDTO.PID, recDTO.Pc)

	// 3. Respondo a CPU
	w.WriteHeader(http.StatusNoContent)
}
