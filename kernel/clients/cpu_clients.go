package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
	"github.com/sisoputnfrba/tp-golang/kernel/models/dtos"
	"net/http"
)

func SendToExec(pcb *models.Pcb, cpu *models.CpuModule) *http.Response {
	config.Logger.Debug(fmt.Sprintf("Sending to Exec PID %d to cpu %s", pcb.Pid, cpu.Id))
	execProcess := dtos.ExecProcessDTO{pcb.Pid, pcb.Pc}
	jsonMessage, err := json.Marshal(execProcess)
	if err != nil {
		config.Logger.Error("ERROR CONVERTING TO JSON")
		panic(err)
	}

	response, err := http.Post(fmt.Sprintf("http://%s/cpu/process/new", cpu.Address), "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		config.Logger.Error("ERROR POSTING")
		panic(err)
	}
	return response
}

func Interrupt(cpu *models.CpuModule) dtos.ReceiveInterrupt {
	sendInterrupt := dtos.SendInterrupt{cpu.Pcb.Pid, "DESALOJO"}
	jsonMessage, err := json.Marshal(sendInterrupt)
	if err != nil {
		config.Logger.Error("ERROR CONVERTING TO JSON")
		panic(err)
	}

	config.Logger.Debug(fmt.Sprintf("Sending INTERRUPT PID %d to cpu %s", cpu.Pcb.Pid, cpu.Id))
	// tengo que esperar a que hagan cpu
	response, err := http.Post(fmt.Sprintf("http://%s/kernel/cpu/interrupt", cpu.Address), "application/json", bytes.NewBuffer(jsonMessage))
	
	if err != nil {
		config.Logger.Error("ERROR POSTING")
		panic(err)
	}
	config.Logger.Debug(fmt.Sprintf("prueba en kernel"))

	var receiveInterrupt dtos.ReceiveInterrupt
	err = json.NewDecoder(response.Body).Decode(&receiveInterrupt)
	if err != nil {
		config.Logger.Error("ERROR DECODING JSON")
		panic(err)
	}
	config.Logger.Debug(fmt.Sprintf("INTERRUPT FINISHED PID %d to cpu %s", receiveInterrupt.PID, cpu.Id))
	return receiveInterrupt
}
