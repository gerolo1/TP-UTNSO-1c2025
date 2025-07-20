package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/io/config"
	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/io/models/dtos"
	"github.com/sisoputnfrba/tp-golang/io/utils"
)

var IoShutdown bool = false
var ShutdownChan = make(chan struct{})

func HandleIORequest(w http.ResponseWriter, r *http.Request) {
	var req dtos.IORequestDTO

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	go func() {
		config.Logger.Debug(fmt.Sprintf("Recibida petición de IO para PID %d - Duración: %d", req.PID, req.Duration))
		config.Logger.Info("## PID: " + strconv.Itoa(req.PID) + " - Inicio de IO - Tiempo: " + strconv.Itoa(req.Duration))

		select {
		case <-time.After(time.Millisecond * time.Duration(req.Duration)):
		case <-ShutdownChan:
		}

		if !IoShutdown {
			err = NotifyKernelIOCompleted(req.PID)
			if err != nil {
				config.Logger.Error("Fallo al notificar al Kernel: " + err.Error())
			}

			config.Logger.Info(fmt.Sprintf("Finalizada petición de IO para PID %d", req.PID))
		} else {
			config.Logger.Info(fmt.Sprintf("No se envia Complete de %d debido a IoShutdown", req.PID))
		}
	}()
}

func NotifyKernelIOCompleted(pid int) error {
	complete := dtos.IOCompleteDTO{
		NameIO: globals.DeviceName,
		PID:    pid,
	}

	config.Logger.Info("## PID: " + strconv.Itoa(pid) + " - Fin de IO - Dispositivo: " + globals.DeviceName)
	return utils.SendToModule(config.Config.IPKernel, config.Config.PortKernel, "kernel/io/complete", complete, nil)
}
