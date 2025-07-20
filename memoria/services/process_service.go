package services

import (
	"errors"
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"github.com/sisoputnfrba/tp-golang/memoria/models/dtos"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"log/slog"
	"strconv"
	"sync"
)

type ProcessService struct {
	logger              *slog.Logger
	fileService         *FileService
	memoryService       *MemoryService
	swapService         *SwapService
	processUtilsService *ProcessUtilsService
	config              *models.Config
}

type ProcessParams struct {
	Logger              *slog.Logger
	FileService         *FileService
	MemoryService       *MemoryService
	SwapService         *SwapService
	ProcessUtilsService *ProcessUtilsService
	Config              *models.Config
}

func NewProcessService(p ProcessParams) *ProcessService {
	return &ProcessService{logger: p.Logger,
		fileService:         p.FileService,
		memoryService:       p.MemoryService,
		swapService:         p.SwapService,
		processUtilsService: p.ProcessUtilsService,
		config:              p.Config}
}

func (s *ProcessService) NewProcess(data dtos.NewProcessDTO) error {
	s.logger.Debug("INIT: ProcessService - NewProcess")

	err := s.memoryService.ChangeValueFreeMemory(false, true, data.Size)
	if err != nil {
		s.logger.Debug("END: ProcessService - NewProcess")
		return err
	}

	newProcess := &models.Process{
		Size:                  data.Size,
		PagesAccessed:         0,
		InstructionsRequested: 0,
		TimesSuspended:        0,
		TimesResumed:          1, // Es 1 ya que representa la cantidad de subidas a memoria principal
		TotalWrites:           0,
		TotalReads:            0,
		Mutex:                 &sync.Mutex{},
	}

	err = s.processUtilsService.AddMetrics(strconv.Itoa(data.PID), newProcess)
	if err != nil {
		s.logger.Debug("END: ProcessService - NewProcess")
		return err
	}

	page0 := &models.Page{
		ID:      0,
		Level:   0,
		Entries: make([]*models.Entry, s.config.EntriesPerPage),
	}
	newPages := &models.ProcessPages{
		Pages:            &[]*models.Page{page0},
		Mutex:            &sync.Mutex{},
		ContadorEntradas: 0,
	}

	s.processUtilsService.AddPages(strconv.Itoa(data.PID), newPages)

	instructionsSlice := s.fileService.OpenInstructionsFile(data.File)

	instructions := &models.ProcessInstructions{
		Instructions: instructionsSlice,
		Mutex:        &sync.Mutex{},
	}

	s.processUtilsService.AddInstructions(strconv.Itoa(data.PID), instructions)

	utils.RequestDelay(s.logger, s.config.MemoryDelay)

	s.logger.Info("## PID: " + strconv.Itoa(data.PID) +
		" - Proceso Creado - Tamaño: " + strconv.Itoa(data.Size))

	s.logger.Debug("END: ProcessService - NewProcess")
	return nil
}

func (s *ProcessService) GetInstruction(pidString string, pcString string) (map[string]string, error) {
	s.logger.Debug("INIT: ProcessService - GetInstruction / PID: " + pidString + " PC: " + pcString)

	pc, _ := strconv.Atoi(pcString)

	instructions, err := s.processUtilsService.FindInstructions(pidString)
	if err != nil {
		s.logger.Debug("END: ProcessService - GetInstruction / PID: " + pidString + " PC: " + pcString)
		return nil, err
	}

	if pc >= len(instructions.Instructions) {
		s.processUtilsService.OkInstructions(instructions)

		s.logger.Error("PC invalid")
		s.logger.Debug("END: ProcessService - GetInstruction / PID: " + pidString + " PC: " + pcString)
		return nil, errors.New("pc invalid")
	}

	instruction := instructions.Instructions[pc]

	s.processUtilsService.OkInstructions(instructions)

	process, err := s.processUtilsService.FindMetrics(pidString)
	if err != nil {
		s.logger.Debug("END: ProcessService - GetInstruction / PID: " + pidString + " PC: " + pcString)
		return nil, err
	}

	process.InstructionsRequested++

	s.processUtilsService.OkMetrics(process)

	utils.RequestDelay(s.logger, s.config.MemoryDelay)

	s.logger.Info("## PID: " + pidString +
		" - Obtener instrucción: " + pcString +
		" - Instrucción: " + instruction)

	response := map[string]string{
		"instruction": instruction,
	}

	s.logger.Debug("END: ProcessService - GetInstruction / PID: " + pidString + " PC: " + pcString)
	return response, nil
}

func (s *ProcessService) SuspendProcess(pidString string) error {
	s.logger.Debug("INIT: ProcessService - SuspendProcess / PID: " + pidString)

	pages, err := s.processUtilsService.FindPages(pidString)
	if err != nil {
		s.logger.Debug("END: ProcessService - SuspendProcess / PID: " + pidString)
		return err
	}

	for _, page := range *pages.Pages {
		if page.Level+1 == s.config.NumberOfLevels {
			for _, entry := range page.Entries {
				if entry == nil {
					break
				}
				data := s.memoryService.FreeMemory(entry.Value, false)

				err := s.swapService.SwapPage(pidString, entry.ID, data)
				if err != nil {
					s.processUtilsService.OkPages(pages)

					s.logger.Debug("END: ProcessService - SuspendProcess / PID: " + pidString)
					return err
				}

				entry.Value = -1

				utils.RequestDelay(s.logger, s.config.SwapDelay)
			}
		}
	}

	s.processUtilsService.OkPages(pages)

	process, err := s.processUtilsService.FindMetrics(pidString)
	if err != nil {
		s.logger.Debug("END: ProcessService - SuspendProcess / PID: " + pidString)
		return err
	}

	_ = s.memoryService.ChangeValueFreeMemory(true, false, process.Size)
	process.TimesSuspended++

	s.processUtilsService.OkMetrics(process)

	s.logger.Debug("END: ProcessService - SuspendProcess / PID: " + pidString)
	return nil
}

func (s *ProcessService) ResumeProcess(pidString string) error {
	s.logger.Debug("INIT: ProcessService - ResumeProcess / PID: " + pidString)

	process, err := s.processUtilsService.FindMetrics(pidString)
	if err != nil {
		s.logger.Debug("END: ProcessService - ResumeProcess / PID: " + pidString)
		return err
	}

	err = s.memoryService.ChangeValueFreeMemory(false, true, process.Size)
	if err != nil {
		s.processUtilsService.OkMetrics(process)

		s.logger.Debug("END: ProcessService - ResumeProcess / PID: " + pidString)
		return err
	}
	process.TimesResumed++

	s.processUtilsService.OkMetrics(process)

	pages, err := s.processUtilsService.FindPages(pidString)
	if err != nil {
		s.logger.Debug("END: ProcessService - ResumeProcess / PID: " + pidString)
		return err
	}

	for _, page := range *pages.Pages {
		if page.Level+1 == s.config.NumberOfLevels {
			for _, entry := range page.Entries {
				if entry == nil {
					break
				}
				data := s.swapService.UnSwapPage(pidString, entry.ID, false)

				newFrame, err := s.memoryService.FindNextFreeFrame()
				if err != nil {
					s.processUtilsService.OkPages(pages)

					s.logger.Debug("END: ProcessService - ResumeProcess / PID: " + pidString)
					return err
				}

				entry.Value = newFrame

				s.memoryService.RestoreUserMemory(data, newFrame)

				utils.RequestDelay(s.logger, s.config.SwapDelay)
			}
		}
	}

	s.processUtilsService.OkPages(pages)

	s.logger.Debug("END: ProcessService - ResumeProcess / PID: " + pidString)
	return nil
}

func (s *ProcessService) ExitProcess(pidString string) error {
	s.logger.Debug("INIT: ProcessService - ExitProcess / PID: " + pidString)

	process := s.processUtilsService.DeleteMetrics(pidString)

	logCloseInfo := "## PID: " + pidString +
		" - Proceso Destruido - Métricas - Acc.T.Pag: " + strconv.Itoa(process.PagesAccessed) +
		"; Inst.Sol.: " + strconv.Itoa(process.InstructionsRequested) +
		"; SWAP: " + strconv.Itoa(process.TimesSuspended) +
		"; Mem.Prin.: " + strconv.Itoa(process.TimesResumed) +
		"; Lec.Mem.: " + strconv.Itoa(process.TotalReads) +
		"; Esc.Mem.: " + strconv.Itoa(process.TotalWrites)

	swaped := process.TimesSuspended == process.TimesResumed

	_ = s.memoryService.ChangeValueFreeMemory(true, false, process.Size)

	pages := s.processUtilsService.DeletePages(pidString)

	for _, page := range *pages.Pages {
		if page.Level+1 == s.config.NumberOfLevels {
			for _, entry := range page.Entries {
				if entry == nil {
					break
				}

				// Check si el proceso esta suspendido o en memoria
				if swaped {
					s.swapService.UnSwapPage(pidString, entry.ID, true)
					utils.RequestDelay(s.logger, s.config.SwapDelay)
				} else {
					_ = s.memoryService.FreeMemory(entry.Value, true)
					utils.RequestDelay(s.logger, s.config.MemoryDelay)
				}
			}
		}
	}

	s.processUtilsService.DeleteInstructions(pidString)

	s.logger.Info(logCloseInfo)

	s.logger.Debug("END: ProcessService - ExitProcess / PID: " + pidString)
	return nil
}
