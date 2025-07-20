package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/memoria/models/dtos"
	"github.com/sisoputnfrba/tp-golang/memoria/services"
)

type KernelController struct {
	logger         *slog.Logger
	memoryService  *services.MemoryService
	processService *services.ProcessService
}

func NewKernelController(logger *slog.Logger, memoryService *services.MemoryService, processService *services.ProcessService) *KernelController {
	return &KernelController{logger: logger, memoryService: memoryService, processService: processService}
}

func (c *KernelController) CreateProcessHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: KernelController - CreateProcessHandler")

	var body dtos.NewProcessDTO

	if request.Body == nil {
		c.logger.Error("BAD REQUEST: BODY NOT FOUND")
		http.Error(response, "BAD REQUEST: BODY NOT FOUND", http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - CreateProcessHandler")
		return
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		c.logger.Error("BAD REQUEST: Invalid JSON format")
		http.Error(response, "BAD REQUEST: Invalid JSON format", http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - CreateProcessHandler")
		return
	}

	err = c.processService.NewProcess(body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - CreateProcessHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)

	c.logger.Debug("~~ END: KernelController - CreateProcessHandler")
}

func (c *KernelController) SuspendProcessHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: KernelController - SuspendProcessHandler")

	pid := request.URL.Query().Get("pid")

	if pid == "" {
		c.logger.Error("BAD QUERY PARAM: PID")
		http.Error(response, "BAD QUERY PARAM AT PID", http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - SuspendProcessHandler")
		return
	}

	err := c.processService.SuspendProcess(pid)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - SuspendProcessHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	c.logger.Debug("~~ END: KernelController - SuspendProcessHandler")
}

func (c *KernelController) ResumeProcessHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: KernelController - ResumeProcessHandler")

	pid := request.URL.Query().Get("pid")

	if pid == "" {
		c.logger.Error("BAD QUERY PARAM: PID")
		http.Error(response, "BAD QUERY PARAM AT PID", http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - ResumeProcessHandler")
		return
	}

	err := c.processService.ResumeProcess(pid)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - ResumeProcessHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	c.logger.Debug("~~ END: KernelController - ResumeProcessHandler")
}

func (c *KernelController) ExitProcessHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: KernelController - ExitProcessHandler")

	pid := request.URL.Query().Get("pid")

	if pid == "" {
		c.logger.Error("BAD QUERY PARAM: PID")
		http.Error(response, "BAD QUERY PARAM AT PID", http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - ExitProcessHandler")
		return
	}

	err := c.processService.ExitProcess(pid)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - ExitProcessHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	c.logger.Debug("~~ END: KernelController - ExitProcessHandler")
}

func (c *KernelController) MemoryDumpHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: KernelController - MemoryDumpHandler")

	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)

		c.logger.Debug("~~ END: KernelController - MemoryDumpHandler")
		return
	}

	var payload map[string]int
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		http.Error(response, "Invalid JSON body", http.StatusBadRequest)

		c.logger.Debug("~~ END: KernelController - MemoryDumpHandler")
		return
	}

	pidString := strconv.Itoa(payload["pid"])
	err = c.memoryService.DumpMemory(pidString)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)

		c.logger.Debug("~~ END: KernelController - MemoryDumpHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	c.logger.Debug("~~ END: KernelController - MemoryDumpHandler")
}
