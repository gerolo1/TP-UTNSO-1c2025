package main

import (
	"github.com/sisoputnfrba/tp-golang/memoria/config"
	"github.com/sisoputnfrba/tp-golang/memoria/controllers"
	"github.com/sisoputnfrba/tp-golang/memoria/services"
	"net/http"
	"strconv"
)

func main() {
	configData := config.InitConfiguration()

	logger := config.InitLogger("memoria.log", configData.LogLevel)

	logger.Debug("INIT ESTRUCTURAS")

	// Estructuras que representan la memoria (generales)
	userMemory := config.InitUserMemory(configData.MemorySize)
	mutexUserMemory := config.InitMutexUserMemory()
	freeMemory := config.InitFreeMemory(configData.MemorySize)
	mutexFreeMemory := config.InitMutexFreeMemory()
	frameBitMap := config.InitFrameBitMap(configData.MemorySize, configData.PageSize)
	mutexFrameBitMap := config.InitMutexFrameBitMap()

	// Estructuras que representan la info del proceso
	processes := config.InitProcesses()
	mutexProcesses := config.InitMutexProcesses()
	processesPages := config.InitProcessesPages()
	mutexProcessesPages := config.InitMutexProcessesPages()
	processesInstructions := config.InitProcessesInstructions()
	mutexProcessesInstructions := config.InitMutexProcessesInstructions()

	// Estructuras sobre SWAP (generales)
	swapBitMap := config.InitSwapBitMap()
	mutexSwapBitMap := config.InitMutexSwapBitMap()
	swapEntries := config.InitSwapEntries()
	mutexSwapEntries := config.InitMutexSwapEntries()

	logger.Debug("OK")

	logger.Debug("INIT INYECCION DE DEPENDENCIAS")

	fileService := services.NewFileService(logger, configData)

	processUtilsParams := services.ProcessUtilsParams{
		Logger:            logger,
		ProcessMap:        processes,
		MutexProcesses:    mutexProcesses,
		InstructionsMap:   processesInstructions,
		MutexInstructions: mutexProcessesInstructions,
		PagesMap:          processesPages,
		MutexPages:        mutexProcessesPages,
	}
	processUtilsService := services.NewProcessUtilsService(processUtilsParams)

	memoryParams := services.MemoryParams{
		Logger:              logger,
		FileService:         fileService,
		ProcessUtilsService: processUtilsService,
		UserMemory:          userMemory,
		MutexUserMemory:     mutexUserMemory,
		FreeMemory:          freeMemory,
		MutexFreeMemory:     mutexFreeMemory,
		Config:              configData,
		FrameBitMap:         frameBitMap,
		MutexFrameBitMap:    mutexFrameBitMap,
	}
	memoryService := services.NewMemoryService(memoryParams)

	paginationParams := services.PaginationParams{
		Logger:              logger,
		MemoryService:       memoryService,
		ProcessUtilsService: processUtilsService,
		Config:              configData,
	}
	paginationService := services.NewPaginationService(paginationParams)

	swapParams := services.SwapParams{
		Logger:           logger,
		FileService:      fileService,
		SwapBitMap:       swapBitMap,
		SwapEntries:      swapEntries,
		MutexSwapBitMap:  mutexSwapBitMap,
		MutexSwapEntries: mutexSwapEntries,
		Config:           configData,
	}
	swapService := services.NewSwapService(swapParams)

	processParams := services.ProcessParams{
		Logger:              logger,
		FileService:         fileService,
		MemoryService:       memoryService,
		SwapService:         swapService,
		ProcessUtilsService: processUtilsService,
		Config:              configData,
	}
	processService := services.NewProcessService(processParams)

	cpuController := controllers.NewCpuController(logger, processService, memoryService, paginationService, configData)
	kernelController := controllers.NewKernelController(logger, memoryService, processService)

	logger.Debug("OK")

	router := config.CrearRouter(cpuController, kernelController)

	addr := configData.IpMemory + ":" + strconv.Itoa(configData.PortMemory)

	logger.Debug("INIT SERVER - " + addr)

	e := http.ListenAndServe(addr, router)
	if e != nil {
		logger.Error("Couldnt init server")
		return
	}
}
