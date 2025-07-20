package controllers

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"github.com/sisoputnfrba/tp-golang/memoria/models/dtos"
	"github.com/sisoputnfrba/tp-golang/memoria/services"
	"log/slog"
	"net/http"
)

type CpuController struct {
	logger            *slog.Logger
	memoryService     *services.MemoryService
	processService    *services.ProcessService
	paginationService *services.PaginationService
	configData        *models.Config
}

func NewCpuController(logger *slog.Logger,
	processService *services.ProcessService,
	memoryService *services.MemoryService,
	paginationService *services.PaginationService,
	configData *models.Config) *CpuController {
	return &CpuController{logger: logger,
		processService:    processService,
		memoryService:     memoryService,
		paginationService: paginationService,
		configData:        configData}
}

func (c *CpuController) GetInstructionHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: CpuController - GetInstructionHandler")

	pid := request.URL.Query().Get("pid")
	pc := request.URL.Query().Get("pc")

	if pid == "" {
		c.logger.Error("BAD QUERY PARAM: PID")
		http.Error(response, "BAD QUERY PARAM AT PID", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - GetInstructionHandler")
		return
	}

	if pc == "" {
		c.logger.Error("BAD QUERY PARAM: PC")
		http.Error(response, "BAD QUERY PARAM AT PC", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - GetInstructionHandler")
		return
	}

	instruction, err := c.processService.GetInstruction(pid, pc)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - GetInstructionHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err = json.NewEncoder(response).Encode(instruction)
	if err != nil {
		c.logger.Error("Cannot Encode Instructions")
		http.Error(response, "Cannot Encode Instructions", http.StatusInternalServerError)

		c.logger.Debug("~~ END: CpuController - GetInstructionHandler")
		return
	}

	c.logger.Debug("~~ END: CpuController - GetInstructionHandler")
}

func (c *CpuController) GetFrameHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: CpuController - GetFrameHandler")

	var body dtos.FrameRequestDTO

	if request.Body == nil {
		c.logger.Error("BAD REQUEST: BODY NOT FOUND")
		http.Error(response, "BAD REQUEST: BODY NOT FOUND", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - WriteMemoryHandler")
		return
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		c.logger.Error("BAD REQUEST: Invalid JSON format")
		http.Error(response, "BAD REQUEST: Invalid JSON format", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - WriteMemoryHandler")
		return
	}

	frameResponse, err := c.paginationService.GetFrame(body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - GetFrameHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err = json.NewEncoder(response).Encode(frameResponse)
	if err != nil {
		c.logger.Error("Cannot Encode Frame")
		http.Error(response, "Cannot Encode Frame", http.StatusInternalServerError)

		c.logger.Debug("~~ END: CpuController - GetFrameHandler")
		return
	}

	c.logger.Debug("~~ END: CpuController - GetFrameHandler")
}

func (c *CpuController) ReadMemoryHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: CpuController - ReadMemoryHandler")

	var body dtos.ReaderDTO

	if request.Body == nil {
		c.logger.Error("BAD REQUEST: BODY NOT FOUND")
		http.Error(response, "BAD REQUEST: BODY NOT FOUND", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - ReadMemoryHandler")
		return
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		c.logger.Error("BAD REQUEST: Invalid JSON format")
		http.Error(response, "BAD REQUEST: Invalid JSON format", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - ReadMemoryHandler")
		return
	}

	memory, err := c.memoryService.ReadMemory(body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - ReadMemoryHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err = json.NewEncoder(response).Encode(memory)
	if err != nil {
		c.logger.Error("Cannot Encode Memory")
		http.Error(response, "Cannot Encode Memory", http.StatusInternalServerError)

		c.logger.Debug("~~ END: CpuController - ReadMemoryHandler")
		return
	}

	c.logger.Debug("~~ END: CpuController - ReadMemoryHandler")
}

func (c *CpuController) WriteMemoryHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: CpuController - WriteMemoryHandler")

	var body dtos.WriterDTO

	if request.Body == nil {
		c.logger.Error("BAD REQUEST: BODY NOT FOUND")
		http.Error(response, "BAD REQUEST: BODY NOT FOUND", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - WriteMemoryHandler")
		return
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		c.logger.Error("BAD REQUEST: Invalid JSON format")
		http.Error(response, "BAD REQUEST: Invalid JSON format", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - WriteMemoryHandler")
		return
	}

	err = c.memoryService.WriteMemory(body, len(body.Value))
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - WriteMemoryHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	c.logger.Debug("~~ END: CpuController - WriteMemoryHandler")
}

func (c *CpuController) ReadMemoryFullPageHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: CpuController - ReadMemoryFullPageHandler")

	var body dtos.ReaderDTO

	if request.Body == nil {
		c.logger.Error("BAD REQUEST: BODY NOT FOUND")
		http.Error(response, "BAD REQUEST: BODY NOT FOUND", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - ReadMemoryFullPageHandler")
		return
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		c.logger.Error("BAD REQUEST: Invalid JSON format")
		http.Error(response, "BAD REQUEST: Invalid JSON format", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - ReadMemoryFullPageHandler")
		return
	}

	body.Size = c.configData.PageSize
	memory, err := c.memoryService.ReadMemory(body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - ReadMemoryFullPageHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err = json.NewEncoder(response).Encode(memory)
	if err != nil {
		c.logger.Error("Cannot Encode Memory")
		http.Error(response, "Cannot Encode Memory", http.StatusInternalServerError)

		c.logger.Debug("~~ END: CpuController - ReadMemoryFullPageHandler")
		return
	}

	c.logger.Debug("~~ END: CpuController - ReadMemoryFullPageHandler")
}

func (c *CpuController) WriteMemoryFullPageHandler(response http.ResponseWriter, request *http.Request) {
	c.logger.Debug("~~ INIT: CpuController - WriteMemoryFullPageHandler")

	var body dtos.WriterDTO

	if request.Body == nil {
		c.logger.Error("BAD REQUEST: BODY NOT FOUND")
		http.Error(response, "BAD REQUEST: BODY NOT FOUND", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - WriteMemoryFullPageHandler")
		return
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		c.logger.Error("BAD REQUEST: Invalid JSON format")
		http.Error(response, "BAD REQUEST: Invalid JSON format", http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - WriteMemoryFullPageHandler")
		return
	}

	err = c.memoryService.WriteMemory(body, c.configData.PageSize)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)

		c.logger.Debug("~~ END: CpuController - WriteMemoryFullPageHandler")
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	c.logger.Debug("~~ END: CpuController - WriteMemoryFullPageHandler")
}

func (c *CpuController) GetConfigCPU(response http.ResponseWriter, _ *http.Request) {
	c.logger.Debug("~~ INIT: CpuController - GetConfigCPU")

	payload := dtos.ConfigCPUDTO{
		NumberOfLevels: c.configData.NumberOfLevels,
		PageSize:       c.configData.PageSize,
		EntriesPerPage: c.configData.EntriesPerPage,
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err := json.NewEncoder(response).Encode(payload)
	if err != nil {
		c.logger.Error("Cannot Encode Configuration")
		http.Error(response, "Cannot Encode Memory", http.StatusInternalServerError)

		c.logger.Debug("~~ END: CpuController - GetConfigCPU")
		return
	}

	c.logger.Debug("~~ END: CpuController - GetConfigCPU")
}
