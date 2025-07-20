package controllers

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/config"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/models/dtos"
	"github.com/sisoputnfrba/tp-golang/cpu/services"
	"net/http"
	"strconv"
)

// 3.a Recibir pid y pc de Kernel para levantar un cpu
func GetPidAndPcFromKernel(response http.ResponseWriter, request *http.Request) {
	var proc struct {
		Pid int `json:"pid"`
		Pc  int `json:"pc"`
	}

	err := json.NewDecoder(request.Body).Decode(&proc)
	if err != nil {
		http.Error(response, "Invalid PID/PC", http.StatusBadRequest)
		return
	}

	//globals.PID = proc.Pid
	//globals.PC = proc.Pc

	go func() {
		services.ProcessInstruction(proc.Pid, proc.Pc)
		globals.InterruptMutex.Lock()
		if globals.InterruptExist {
			globals.InterruptExist = false
			globals.InterruptSync.Unlock()
			config.Logger.Debug("INTERRUPTSYNC UNLOCK EN GETPIDANDPCFROMKERNEL")
		}
		globals.InterruptMutex.Unlock()
	}()

	response.WriteHeader(http.StatusOK)
}

func ReceiveInterrupt(w http.ResponseWriter, r *http.Request) {
	var interrupt dtos.Interrupt
	//Recibo pedido de intrp
	err := json.NewDecoder(r.Body).Decode(&interrupt)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		config.Logger.Error("Error al decodificar interrupción: " + err.Error())
		return
	}
	config.Logger.Info("## Llega interrupción al puerto Interrupt - PID: " + strconv.Itoa(interrupt.PID) + " - Motivo: " + interrupt.Motivo)

	globals.InterruptMutex.Lock()
	config.Logger.Debug("INTERRUPTMUTEX LOCK EN RECEIVE INTERRUPT")

	// ponemos el bool en true
	globals.InterruptExist = true
	// ponemos la estructura
	globals.Interrupt = interrupt

	// unlockeamos
	globals.InterruptMutex.Unlock()
	config.Logger.Debug("INTERRUPTMUTEX UNLOCK EN RECEIVE INTERRUPT")

	// esperamos (lock)
	globals.InterruptSync.Lock()
	config.Logger.Debug("INTERRUPTSYNC LOCK EN RECEIVE INTERRUPT")
	config.Logger.Debug("esta es la prueba")

	// escribimos respuesta
	processInterrupt := dtos.ProcessInterrupt{
		PID:    interrupt.PID,
		PC:     globals.PC,
		Motivo: interrupt.Motivo,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(processInterrupt)

	config.Logger.Debug("Interrupción atendida para PID " + strconv.Itoa(interrupt.PID))
}
