package services

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/io/config"
	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/io/models/dtos"
	"github.com/sisoputnfrba/tp-golang/io/utils"
)

func SendIOHandshakeToKernel() {
	payload := dtos.IOHandshakeDTO{
		Name: globals.DeviceName,
		IP:   config.Config.IPIO,
		Port: config.Config.PortIO,
	}

	var response dtos.ResponseDTO
	err := utils.SendToModule(config.Config.IPKernel, config.Config.PortKernel, "kernel/io/handshake", payload, &response)
	config.Logger.Info("Respuesta del Kernel: " + response.Message)

	if err != nil {
		config.Logger.Error("Error al enviar handshake al kernel: " + err.Error())
		os.Exit(1)
	}
	config.Logger.Info("Handshake con Kernel exitoso para IO: " + globals.DeviceName)
}

func SendIOShutdownToKernel() {
	IoShutdown = true
	close(ShutdownChan)

	payload := dtos.IOHandshakeDTO{
		Name: globals.DeviceName,
		IP:   config.Config.IPIO,
		Port: config.Config.PortIO,
	}

	err := utils.SendToModule(config.Config.IPKernel, config.Config.PortKernel, "kernel/io/shutdown", payload, nil)
	if err != nil {
		config.Logger.Error(fmt.Sprintf("IO: error al enviar shutdown al Kernel: %v", err))
	} else {
		config.Logger.Info("IO finalizado correctamente.")
	}

}
