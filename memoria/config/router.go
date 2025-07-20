package config

import (
	"github.com/sisoputnfrba/tp-golang/memoria/controllers"
	"net/http"
)

func CrearRouter(cpuController *controllers.CpuController,
	kernelController *controllers.KernelController) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/memory/processes/new", kernelController.CreateProcessHandler)
	mux.HandleFunc("/memory/processes/suspend", kernelController.SuspendProcessHandler)
	mux.HandleFunc("/memory/processes/resume", kernelController.ResumeProcessHandler)
	mux.HandleFunc("/memory/processes/exit", kernelController.ExitProcessHandler)
	mux.HandleFunc("/memory/dump", kernelController.MemoryDumpHandler)

	mux.HandleFunc("/memory/instruction", cpuController.GetInstructionHandler)
	mux.HandleFunc("/memory/frame", cpuController.GetFrameHandler)
	mux.HandleFunc("/memory/write", cpuController.WriteMemoryHandler)
	mux.HandleFunc("/memory/read", cpuController.ReadMemoryHandler)
	mux.HandleFunc("/memory/full-write", cpuController.WriteMemoryFullPageHandler)
	mux.HandleFunc("/memory/full-read", cpuController.ReadMemoryFullPageHandler)
	mux.HandleFunc("/memory/config", cpuController.GetConfigCPU)

	return mux
}
