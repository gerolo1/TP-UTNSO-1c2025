package main

import (
	"github.com/sisoputnfrba/tp-golang/cpu/config"
	"github.com/sisoputnfrba/tp-golang/cpu/controllers"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/models"
	"github.com/sisoputnfrba/tp-golang/cpu/services"
	"net/http"
	"os"
	"sync"
	"strconv"
)

func main() {

	identifier := "1"
	var port string

	if len(os.Args) == 3 {
		identifier = os.Args[1]
		port = os.Args[2]
	}

	//Config identificador
	globals.ConfigFile = config.InitConfiguration()
	config.InitLogger("cpu_" + identifier + ".log")

	config.Logger.Debug("Identificador CPU: " + identifier)

	if port != "" {
		globals.ConfigFile.PortCPU, _ = strconv.Atoi(port)
	}

	initGlobals(identifier)

	//Esperar mensaje de Kernel
	mux := http.NewServeMux()

	//Handlers de CPU
	mux.HandleFunc("/cpu/process/new", controllers.GetPidAndPcFromKernel)
	mux.HandleFunc("/kernel/cpu/interrupt", controllers.ReceiveInterrupt)

	//Direccion cpu
	addr := globals.ConfigFile.IPCPU + ":" + strconv.Itoa(globals.ConfigFile.PortCPU)

	config.Logger.Debug("START SERVER on " + addr)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		config.Logger.Error("Couldnt init CPU server: " + err.Error())
		return
	}
}

func initGlobals(identifier string) {
	//Configurar TLB
	globals.Tlb = &models.TLB{
		Entries:    make([]models.TLBEntry, 0, globals.ConfigFile.TLBEntries),
		MaxEntries: globals.ConfigFile.TLBEntries,
		OrderCont:  0,
		AccesCont:  0,
	}

	//Configurar Cache
	globals.Cache = &models.Cache{
		Entries:      make([]models.CacheEntry, 0, globals.ConfigFile.CacheEntries),
		MaxEntries:   globals.ConfigFile.CacheEntries,
		ClockPointer: 0,
	}

	// Handshake con el kernel
	services.DoHandshake(globals.ConfigFile.IPKernel, globals.ConfigFile.PortKernel, globals.ConfigFile.IPCPU, globals.ConfigFile.PortCPU, identifier)

	// Obtener config de paginación multinivel desde Memoria (tamano de pagina, niveles, entradas por tabla))
	var errConfig error
	globals.PageSize, globals.Levels, globals.EntriesPerTable, errConfig = services.GetPageTableConfig(globals.ConfigFile.IPMemory, globals.ConfigFile.PortMemory)
	if errConfig != nil {
		config.Logger.Error("Error al obtener config de tablas de páginas: " + errConfig.Error())
	}

	globals.Identifier = identifier

	globals.InterruptSync = &sync.Mutex{}
	globals.InterruptSync.Lock() //
}
