package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/controllers"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
	"github.com/sisoputnfrba/tp-golang/kernel/services"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
)

func main() {
	nameFirstProcess := config.InitConfiguration()
	sizeFirstProcess := config.FirstProcessSize(nameFirstProcess)
	config.InitLogger("kernel.log")

	mux := http.NewServeMux()

	// devuelve 202 (statusOK)
	mux.HandleFunc("/kernel/cpu/new", controllers.NewCpuHandler)
	// devuelve 204 (http.StatusNoContent)
	mux.HandleFunc("/kernel/cpu/exit", controllers.ExitHandler)
	mux.HandleFunc("/kernel/cpu/dump", controllers.DumpHandler)
	mux.HandleFunc("/kernel/cpu/init", controllers.InitProcessHandler)
	mux.HandleFunc("/kernel/cpu/io", controllers.IoHandler)

	mux.HandleFunc("/kernel/io/complete", controllers.IoCompleteHandler)
	mux.HandleFunc("/kernel/io/handshake", controllers.IoHandshakeHandler)
	mux.HandleFunc("/kernel/io/shutdown", controllers.IoShutdownHandler)

	addr := config.Config.IpKernel + ":" + strconv.Itoa(config.Config.PortKernel)

	config.Logger.Debug("START SERVER")

	go func() {
		fmt.Println("Presiona ENTER para continuar")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Println("Continuando")

		pcbKernel := &models.Pcb{0, 0, sizeFirstProcess, nameFirstProcess, services.EmptyStateMap(), "", 0, 0, config.Config.InitialEstimate}
		services.NewSem.Lock()
		utils.LockLogNew(pcbKernel.Pid, "MAIN")
		services.ToNew(pcbKernel)
	}()

	e := http.ListenAndServe(addr, mux)
	if e != nil {
		config.Logger.Error("Couldnt init server")
		return
	}
}
