/*

package controllers
import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/config"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/models/dtos"
)

func HOLA(w http.ResponseWriter, r *http.Request) {
	var interrupt dtos.Interrupt
	//Recibo pedido de intrp
	err := json.NewDecoder(r.Body).Decode(&interrupt)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		config.Logger.Error("Error al decodificar interrupción: " + err.Error())
		return
	}
	config.Logger.Info("## Llega interrupción al puerto Interrupt - PID: " + strconv.Itoa(globals.Interrupt.PID) + " - Motivo: " + globals.Interrupt.Motivo)

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

	// escribimos respuesta
	processInterrupt := dtos.ProcessInterrupt{
		PID:    interrupt.PID,
		PC:     globals.PC,
		Motivo: interrupt.Motivo,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(processInterrupt)

	config.Logger.Info("Interrupción atendida para PID " + strconv.Itoa(interrupt.PID))
}
*/