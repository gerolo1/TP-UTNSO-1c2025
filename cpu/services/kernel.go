package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"io"

	"github.com/sisoputnfrba/tp-golang/cpu/config"
	"github.com/sisoputnfrba/tp-golang/cpu/models/dtos"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
)

func DoHandshake(targetIP string, targetPort int, myIp string, myPort int, myId string) {
	handshake := dtos.Handshake{
		Ip:   myIp,
		Port: myPort,
		Id:   myId,
	}

	statusCode, body, err := utils.SendToModuleWithStatus(targetIP, targetPort, "kernel/cpu/new", handshake)
	if err != nil {
		config.Logger.Error("Error al enviar handshake: %v\n", err)
		return
	}

	if statusCode != http.StatusOK {
		config.Logger.Error("Kernel respondió con error. Código: %d, Mensaje: %s\n", statusCode, string(body))
		return
	}

	config.Logger.Info("Handshake con kernel exitoso.")
}


func SendDumpMemoryToKernel(pid int, pc int, targetIP string, targetPort int) (string, error) {
	reqDTO := dtos.DumpRequestDTO{
		PID: pid,
		Pc:  pc,
	}

	// Serializar JSON
	jsonData, err := json.Marshal(reqDTO)
	if err != nil {
		return "", fmt.Errorf("error serializing request: %v", err)
	}

	// Armar URL
	url := fmt.Sprintf("http://%s:%d/kernel/cpu/dump", targetIP, targetPort)

	// Hacer POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error contacting kernel: %v", err)
	}
	defer resp.Body.Close()


	// Verificar que la respuesta sea 204 No Content
	if resp.StatusCode != http.StatusNoContent {
		// Podés leer el cuerpo en caso de error para obtener más información
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	// Si llegó hasta acá, fue exitosa con 204 No Content
	return "", nil
}

func SendExitToKernel(pid int, pc int, targetIP string, targetPort int) error {
	proc := dtos.ExitProcess{
		Pid: pid,
		Pc:  pc,
	}
	return utils.SendToModule(targetIP, targetPort, "kernel/cpu/exit", proc, nil)

}


func SendIOToKernel(pid int, pc int, instruction string, targetIP string, targetPort int) error {
	parts := strings.Split(instruction, " ")

	device := parts[1]
	duration, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}

	ioReq := dtos.IORequest{
		PID:      pid,
		Device:   device,
		Duration: duration,
		PC:       pc,
	}

	return utils.SendToModule(targetIP, targetPort, "kernel/cpu/io", ioReq, nil)
}

func SendInitProcToKernel(instruction string, targetIP string, targetPort int) error {
	parts := strings.Split(instruction, " ")
	filename := parts[1]
	size, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}

	initReq := dtos.InitProcRequest{
		Filename: filename,
		Size:     size,
	}

	return utils.SendToModule(targetIP, targetPort, "kernel/cpu/init", initReq, nil)
}

func SendInterruptToKernel(interrupt dtos.ProcessInterrupt, ip string, port int) error {
	var respuesta interface{} // no espero nada
	err := utils.SendToModule(ip, port, "/cpu/interrupt", interrupt, &respuesta)
	if err != nil {
		config.Logger.Error("Error al enviar interrupción al kernel: " + err.Error())
	}

	
	return err
}
