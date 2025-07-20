package services

import (
	"errors"
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"github.com/sisoputnfrba/tp-golang/memoria/models/dtos"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

type MemoryService struct {
	logger              *slog.Logger
	fileService         *FileService
	processUtilsService *ProcessUtilsService
	config              *models.Config
	userMemory          []byte
	mutexUserMemory     *sync.Mutex
	freeMemory          *int
	mutexFreeMemory     *sync.Mutex
	frameBitMap         []bool
	mutexFrameBitMap    *sync.Mutex
}

type MemoryParams struct {
	Logger              *slog.Logger
	FileService         *FileService
	ProcessUtilsService *ProcessUtilsService
	Config              *models.Config
	UserMemory          []byte
	MutexUserMemory     *sync.Mutex
	FreeMemory          *int
	MutexFreeMemory     *sync.Mutex
	FrameBitMap         []bool
	MutexFrameBitMap    *sync.Mutex
}

func NewMemoryService(p MemoryParams) *MemoryService {
	return &MemoryService{logger: p.Logger,
		fileService:         p.FileService,
		processUtilsService: p.ProcessUtilsService,
		userMemory:          p.UserMemory,
		mutexUserMemory:     p.MutexUserMemory,
		freeMemory:          p.FreeMemory,
		mutexFreeMemory:     p.MutexFreeMemory,
		frameBitMap:         p.FrameBitMap,
		mutexFrameBitMap:    p.MutexFrameBitMap,
		config:              p.Config}
}

func (s *MemoryService) FindNextFreeFrame() (int, error) {
	s.logger.Debug("INIT: MemoryService - FindNextFreeFrame")

	s.mutexFrameBitMap.Lock()
	for i, val := range s.frameBitMap {
		if !val {
			s.frameBitMap[i] = true

			s.mutexFrameBitMap.Unlock()
			return i * s.config.PageSize, nil
		}
	}

	s.mutexFrameBitMap.Unlock()

	s.logger.Debug("END: MemoryService - FindNextFreeFrame")
	return 0, errors.New("FATAL: bitmap full")
}

func (s *MemoryService) FreeMemory(frame int, exiting bool) []byte {
	s.logger.Debug("INIT: MemoryService - FreeMemory")

	s.mutexFrameBitMap.Lock()
	s.frameBitMap[frame/s.config.PageSize] = false
	s.mutexFrameBitMap.Unlock()

	if exiting {
		s.logger.Debug("END: MemoryService - FreeMemory")
		return nil
	}

	data := make([]byte, s.config.PageSize)

	s.mutexUserMemory.Lock()
	copy(data, s.userMemory[frame:frame+s.config.PageSize])
	s.mutexUserMemory.Unlock()

	s.logger.Debug("END: MemoryService - FreeMemory")
	return data

}

func (s *MemoryService) RestoreUserMemory(data []byte, dir int) {
	s.logger.Debug("INIT: MemoryService - RestoreUserMemory")

	s.mutexUserMemory.Lock()
	copy(s.userMemory[dir:dir+s.config.PageSize], data)
	s.mutexUserMemory.Unlock()

	s.logger.Debug("END: MemoryService - RestoreUserMemory")
}

func (s *MemoryService) ReadMemory(data dtos.ReaderDTO) (map[string][]byte, error) {
	s.logger.Debug("INIT: MemoryService - ReadMemory / PID: " + strconv.Itoa(data.PID))

	s.mutexUserMemory.Lock()
	if data.DF < 0 || data.DF > len(s.userMemory) {
		s.mutexUserMemory.Unlock()
		s.logger.Error("Invalid access to memory")

		s.logger.Debug("END: MemoryService - ReadMemory / PID: " + strconv.Itoa(data.PID))
		return nil, errors.New("invalid access to memory")
	}

	responseData := make([]byte, data.Size)

	copy(responseData, s.userMemory[data.DF:data.DF+data.Size])
	s.mutexUserMemory.Unlock()

	process, err := s.processUtilsService.FindMetrics(strconv.Itoa(data.PID))
	if err != nil {
		s.logger.Debug("END: MemoryService - ReadMemory / PID: " + strconv.Itoa(data.PID))
		return nil, err
	}

	process.TotalReads++

	s.processUtilsService.OkMetrics(process)

	utils.RequestDelay(s.logger, s.config.MemoryDelay)

	s.logger.Info("## PID: " + strconv.Itoa(data.PID) +
		" - Lectura - Dir. Física: " + strconv.Itoa(data.DF) +
		" - Tamaño: " + strconv.Itoa(data.Size))

	response := map[string][]byte{
		"data": responseData,
	}

	s.logger.Debug("END: MemoryService - ReadMemory / PID: " + strconv.Itoa(data.PID))
	return response, nil
}

func (s *MemoryService) WriteMemory(data dtos.WriterDTO, size int) error {
	s.logger.Debug("INIT: MemoryService - ReadMemory / PID: " + strconv.Itoa(data.PID))

	s.mutexUserMemory.Lock()
	if data.DF < 0 ||
		data.DF > len(s.userMemory) ||
		data.DF/s.config.PageSize != (data.DF+size)/s.config.PageSize {
		s.mutexUserMemory.Unlock()
		s.logger.Error("Invalid access to memory")

		s.logger.Debug("END: MemoryService - WriteMemory / PID: " + strconv.Itoa(data.PID))
		return errors.New("invalid access to memory")
	}

	copy(s.userMemory[data.DF:data.DF+size], data.Value)
	s.mutexUserMemory.Unlock()

	process, err := s.processUtilsService.FindMetrics(strconv.Itoa(data.PID))
	if err != nil {
		s.logger.Debug("END: MemoryService - WriteMemory / PID: " + strconv.Itoa(data.PID))
		return err
	}

	process.TotalWrites++

	s.processUtilsService.OkMetrics(process)

	utils.RequestDelay(s.logger, s.config.MemoryDelay)

	s.logger.Info("## PID: " + strconv.Itoa(data.PID) +
		" - Escritura - Dir. Física: " + strconv.Itoa(data.DF) +
		" - Tamaño: " + strconv.Itoa(size))

	s.logger.Debug("END: MemoryService - WriteMemory / PID: " + strconv.Itoa(data.PID))
	return nil
}

func (s *MemoryService) DumpMemory(pidString string) error {
	s.logger.Debug("INIT: MemoryService - DumpMemory / PID: " + pidString)

	s.logger.Info("## PID: " + pidString + " - Memory Dump solicitado")

	pages, err := s.processUtilsService.FindPages(pidString)
	if err != nil {
		s.logger.Debug("END: MemoryService - DumpMemory - PID: " + pidString)
		return err
	}

	timestamp := time.Now()
	hora := timestamp.Format("15:04:05")
	fileName := pidString + "-" + hora + ".dmp"

	var data []byte
	contEntries := 0

	for _, page := range *pages.Pages {
		if page.Level+1 == s.config.NumberOfLevels {
			for _, entry := range page.Entries {
				if entry == nil {
					break
				}
				dataEntry := make([]byte, s.config.PageSize)

				s.mutexUserMemory.Lock()
				copy(dataEntry, s.userMemory[entry.Value:entry.Value+s.config.PageSize])
				s.mutexUserMemory.Unlock()

				data = append(data, dataEntry...)
				contEntries++
			}
		}
	}

	s.processUtilsService.OkPages(pages)

	s.logger.Debug("File Name: " + fileName)
	s.fileService.CreateDumpFile(fileName, data)

	utils.RequestDelay(s.logger, s.config.MemoryDelay)

	s.logger.Debug("END: MemoryService - DumpMemory / PID: " + pidString)

	return nil
}

func (s *MemoryService) ChangeValueFreeMemory(isAdd bool, checkSpace bool, value int) error {
	s.mutexFreeMemory.Lock()
	if checkSpace && *s.freeMemory < value {
		s.mutexFreeMemory.Unlock()
		s.logger.Debug("No space to initialize")
		return errors.New("no space to initialize")
	}

	if isAdd {
		s.logger.Debug("FreeMemory - Before: " + strconv.Itoa(*s.freeMemory) +
			" After: " + strconv.Itoa(*s.freeMemory+value))
		*s.freeMemory += value
	} else {
		s.logger.Debug("FreeMemory - Before: " + strconv.Itoa(*s.freeMemory) +
			" After: " + strconv.Itoa(*s.freeMemory-value))
		*s.freeMemory -= value
	}
	s.mutexFreeMemory.Unlock()
	return nil
}
