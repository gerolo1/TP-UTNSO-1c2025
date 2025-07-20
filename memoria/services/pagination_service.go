package services

import (
	"errors"
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"github.com/sisoputnfrba/tp-golang/memoria/models/dtos"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"log/slog"
	"strconv"
)

type PaginationService struct {
	logger              *slog.Logger
	memoryService       *MemoryService
	ProcessUtilsService *ProcessUtilsService
	config              *models.Config
}

type PaginationParams struct {
	Logger              *slog.Logger
	MemoryService       *MemoryService
	ProcessUtilsService *ProcessUtilsService
	Config              *models.Config
}

func NewPaginationService(p PaginationParams) *PaginationService {
	return &PaginationService{logger: p.Logger,
		memoryService:       p.MemoryService,
		ProcessUtilsService: p.ProcessUtilsService,
		config:              p.Config}
}

func (s *PaginationService) GetFrame(data dtos.FrameRequestDTO) (map[string]int, error) {
	s.logger.Debug("INIT: PaginationService - GetFrame / PID: " + strconv.Itoa(data.PID))

	if len(data.Entries) != s.config.NumberOfLevels {
		s.logger.Error("BAD DF")

		s.logger.Debug("END: PaginationService - GetFrame / PID: " + strconv.Itoa(data.PID))
		return nil, errors.New("BAD DF")
	}

	actualPage := 0

	pages, err := s.ProcessUtilsService.FindPages(strconv.Itoa(data.PID))
	if err != nil {
		s.logger.Debug("END: PaginationService - GetFrame / PID: " + strconv.Itoa(data.PID))
		return nil, err
	}

	// Busco todas las entradas en la lista de paginas
	for _, dl := range data.Entries {
		for _, page := range *pages.Pages {

			//Si encuentra la pagina indicada procede a buscar en las entradas
			if page.ID == actualPage {

				//Buscando por si existe la entrada
				for i, entry := range page.Entries {
					if entry == nil {
						if page.Level+1 == s.config.NumberOfLevels {
							var err error
							actualPage, err = s.memoryService.FindNextFreeFrame()
							if err != nil {
								s.ProcessUtilsService.OkPages(pages)

								s.logger.Debug("END: PaginationService - GetFrame / PID: " + strconv.Itoa(data.PID))
								return nil, err
							}
						} else {
							actualPage = len(*pages.Pages)
							newPage := &models.Page{
								ID:      actualPage,
								Level:   page.Level + 1,
								Entries: make([]*models.Entry, s.config.EntriesPerPage),
							}
							*pages.Pages = append(*pages.Pages, newPage)
						}

						newEntry := &models.Entry{
							Key:   dl,
							Value: actualPage,
							ID:    pages.ContadorEntradas,
						}
						pages.ContadorEntradas++

						page.Entries[i] = newEntry
						break
					} else if entry.Key == dl {
						actualPage = entry.Value
						break
					}
				}

				utils.RequestDelay(s.logger, s.config.MemoryDelay)
				break
			}
		}
	}

	s.ProcessUtilsService.OkPages(pages)

	process, err := s.ProcessUtilsService.FindMetrics(strconv.Itoa(data.PID))
	if err != nil {
		s.logger.Debug("END: PaginationService - GetFrame / PID: " + strconv.Itoa(data.PID))
		return nil, err
	}

	process.PagesAccessed += s.config.NumberOfLevels

	s.ProcessUtilsService.OkMetrics(process)

	response := map[string]int{
		"frame": actualPage,
	}

	s.logger.Debug("END: PaginationService - GetFrame / PID: " + strconv.Itoa(data.PID))

	return response, nil
}
