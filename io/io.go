package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sisoputnfrba/tp-golang/io/config"
	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/io/services"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Falta el nombre del dispositivo IO.")
		os.Exit(1)
	}

	globals.DeviceName = os.Args[1]
	var port string = os.Args[2]

	config.InitConfiguration()
	portInt, _ := strconv.Atoi(port)
	config.Config.PortIO = portInt
	config.InitLogger(globals.DeviceName + ".log")

	config.Logger.Info(">>> EJECUTANDO IO DEBUG <<<") //Para saber que lo levanté

	// Handshake con Kernel
	services.SendIOHandshakeToKernel()

	// Levantar servidor
	mux := http.NewServeMux()
	mux.HandleFunc("/io/request", services.HandleIORequest)
	addr := config.Config.IPIO + ":" + strconv.Itoa(config.Config.PortIO)

	config.Logger.Info("Servidor IO corriendo en " + addr) // para saber si llego bien en pruebas

	go http.ListenAndServe(addr, mux)

	// Escuchar señales para finalización
	//creo canal para recibir senales del s.o.
	c := make(chan os.Signal, 1)
	//cuando recibo senales las mando al canal c
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	//se bloquea esperando a q llegue senal
	<-c

	// Avisar a Kernel que se cierra
	services.SendIOShutdownToKernel()
}
