package services

import (
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"log/slog"
	"slices"
	"strconv"
	"sync"
)

type SwapService struct {
	logger           *slog.Logger
	fileService      *FileService
	config           *models.Config
	swapBitMap       []bool
	mutexSwapBitMap  *sync.Mutex
	swapEntries      *[]*models.SwapEntry
	mutexSwapEntries *sync.Mutex
}

type SwapParams struct {
	Logger           *slog.Logger
	FileService      *FileService
	Config           *models.Config
	SwapBitMap       []bool
	MutexSwapBitMap  *sync.Mutex
	SwapEntries      *[]*models.SwapEntry
	MutexSwapEntries *sync.Mutex
}

func NewSwapService(p SwapParams) *SwapService {
	return &SwapService{logger: p.Logger,
		fileService:      p.FileService,
		swapBitMap:       p.SwapBitMap,
		swapEntries:      p.SwapEntries,
		mutexSwapBitMap:  p.MutexSwapBitMap,
		mutexSwapEntries: p.MutexSwapEntries,
		config:           p.Config}
}

func (s *SwapService) SwapPage(pid string, entryID int, data []byte) error {
	s.logger.Debug("INIT: SwapService - SwapPage / PID: " + pid)

	frameSwap := -1

	s.mutexSwapBitMap.Lock()

	for i, val := range s.swapBitMap {
		if !val {
			s.swapBitMap[i] = true
			frameSwap = i * s.config.PageSize
		}
	}

	if frameSwap == -1 {
		s.swapBitMap = append(s.swapBitMap, true)
		frameSwap = (len(s.swapBitMap) - 1) * s.config.PageSize
	}

	s.mutexSwapBitMap.Unlock()

	pidInt, _ := strconv.Atoi(pid)
	newEntry := &models.SwapEntry{
		PID:       pidInt,
		EntryID:   entryID,
		SwapFrame: frameSwap,
	}

	s.mutexSwapEntries.Lock()
	*s.swapEntries = append(*s.swapEntries, newEntry)
	s.mutexSwapEntries.Unlock()

	s.fileService.WriteSwapFile(data, frameSwap)

	s.logger.Debug("END: SwapService - SwapPage / PID: " + pid)
	return nil
}

func (s *SwapService) UnSwapPage(pid string, entryID int, exiting bool) []byte {
	s.logger.Debug("INIT: SwapService - UnSwapPage / PID: " + pid)

	pidInt, _ := strconv.Atoi(pid)

	s.mutexSwapEntries.Lock()
	for i, swapEntry := range *s.swapEntries {
		if swapEntry.PID == pidInt && swapEntry.EntryID == entryID {
			s.FreeFrameSwap(swapEntry.SwapFrame)

			*s.swapEntries = slices.Delete(*s.swapEntries, i, i+1)

			if exiting {
				s.mutexSwapEntries.Unlock()

				s.logger.Debug("END: SwapService - UnSwapPage / PID: " + pid)
				return nil
			}

			data := s.fileService.ReadSwapFile(swapEntry.SwapFrame)

			s.mutexSwapEntries.Unlock()

			s.logger.Debug("END: SwapService - UnSwapPage / PID: " + pid)
			return data
		}
	}

	s.mutexSwapEntries.Unlock()

	s.logger.Debug("END: SwapService - UnSwapPage / PID: " + pid)
	return nil
}

func (s *SwapService) FreeFrameSwap(frame int) {
	s.logger.Debug("INIT: SwapService - FreeFrameSwap")

	s.mutexSwapBitMap.Lock() //ToDo: Siempre que se swapea, se extiende el swap, no se sobreescribe si ya se libero
	s.swapBitMap[frame/s.config.PageSize] = false
	s.mutexSwapBitMap.Unlock()

	s.logger.Debug("END: SwapService - FreeFrameSwap")
}
