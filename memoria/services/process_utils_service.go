package services

import (
	"errors"
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"log/slog"
	"sync"
)

type ProcessUtilsService struct {
	logger            *slog.Logger
	processMap        map[string]*models.Process
	mutexProcesses    *sync.RWMutex
	instructionsMap   map[string]*models.ProcessInstructions
	mutexInstructions *sync.RWMutex
	pagesMap          map[string]*models.ProcessPages
	mutexPages        *sync.RWMutex
}

type ProcessUtilsParams struct {
	Logger            *slog.Logger
	ProcessMap        map[string]*models.Process
	MutexProcesses    *sync.RWMutex
	InstructionsMap   map[string]*models.ProcessInstructions
	MutexInstructions *sync.RWMutex
	PagesMap          map[string]*models.ProcessPages
	MutexPages        *sync.RWMutex
}

func NewProcessUtilsService(p ProcessUtilsParams) *ProcessUtilsService {
	return &ProcessUtilsService{logger: p.Logger,
		processMap:        p.ProcessMap,
		mutexProcesses:    p.MutexProcesses,
		instructionsMap:   p.InstructionsMap,
		mutexInstructions: p.MutexInstructions,
		pagesMap:          p.PagesMap,
		mutexPages:        p.MutexPages}
}

func (s *ProcessUtilsService) FindInstructions(pid string) (*models.ProcessInstructions, error) {
	s.mutexInstructions.RLock()

	instructions, exists := s.instructionsMap[pid]
	if !exists {
		s.mutexInstructions.RUnlock()

		s.logger.Error("Instructions not exists")
		return nil, errors.New("instructions not exists")
	}

	instructions.Mutex.Lock()

	return instructions, nil
}

func (s *ProcessUtilsService) OkInstructions(instructions *models.ProcessInstructions) {
	instructions.Mutex.Unlock()
	s.mutexInstructions.RUnlock()
}

func (s *ProcessUtilsService) FindPages(pid string) (*models.ProcessPages, error) {
	s.mutexPages.RLock()

	pages, exists := s.pagesMap[pid]
	if !exists {
		s.mutexPages.RUnlock()

		s.logger.Error("Pages not exists")
		return nil, errors.New("pages not exists")
	}

	pages.Mutex.Lock()

	return pages, nil
}

func (s *ProcessUtilsService) OkPages(pages *models.ProcessPages) {
	pages.Mutex.Unlock()
	s.mutexPages.RUnlock()
}

func (s *ProcessUtilsService) FindMetrics(pid string) (*models.Process, error) {
	s.mutexProcesses.RLock()

	processes, exists := s.processMap[pid]
	if !exists {
		s.mutexProcesses.RUnlock()

		s.logger.Error("Process not exists")
		return nil, errors.New("process not exists")
	}

	processes.Mutex.Lock()

	return processes, nil
}

func (s *ProcessUtilsService) OkMetrics(process *models.Process) {
	process.Mutex.Unlock()
	s.mutexProcesses.RUnlock()
}

func (s *ProcessUtilsService) AddMetrics(pid string, newProcess *models.Process) error {
	s.mutexProcesses.Lock()

	_, exists := s.processMap[pid]
	if exists {
		s.mutexProcesses.Unlock()

		s.logger.Error("Process already exists")
		return errors.New("process already exists")
	}

	s.processMap[pid] = newProcess

	s.mutexProcesses.Unlock()

	return nil
}

func (s *ProcessUtilsService) AddPages(pid string, newPages *models.ProcessPages) {
	s.mutexPages.Lock()

	s.pagesMap[pid] = newPages

	s.mutexPages.Unlock()
}

func (s *ProcessUtilsService) AddInstructions(pid string, newInstructions *models.ProcessInstructions) {
	s.mutexInstructions.Lock()

	s.instructionsMap[pid] = newInstructions

	s.mutexInstructions.Unlock()
}

func (s *ProcessUtilsService) DeleteMetrics(pid string) models.Process {
	s.mutexProcesses.Lock()
	process, _ := (s.processMap)[pid]
	delete(s.processMap, pid)
	s.mutexProcesses.Unlock()

	return *process
}

func (s *ProcessUtilsService) DeletePages(pid string) models.ProcessPages {
	s.mutexPages.Lock()
	pages, _ := (s.pagesMap)[pid]
	delete(s.pagesMap, pid)
	s.mutexPages.Unlock()

	return *pages
}

func (s *ProcessUtilsService) DeleteInstructions(pid string) {
	s.mutexInstructions.Lock()
	delete(s.instructionsMap, pid)
	s.mutexInstructions.Unlock()
}
