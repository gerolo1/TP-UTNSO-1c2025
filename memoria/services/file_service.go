package services

import (
	"bufio"
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"log/slog"
	"os"
)

type FileService struct {
	logger *slog.Logger
	config *models.Config
}

func NewFileService(logger *slog.Logger, config *models.Config) *FileService {
	return &FileService{logger: logger, config: config}
}

func (s *FileService) OpenInstructionsFile(fileName string) []string {
	s.logger.Debug("INIT: FileService - OpenInstructionsFile")

	var instructions []string
	file, e := os.Open(s.config.ScriptsPath + fileName)

	if e != nil {
		s.logger.Error("File Problem: Cannot Open File")

		s.logger.Debug("END: FileService - OpenInstructionsFile")
		panic(e)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		instructions = append(instructions, line)
	}

	s.logger.Debug("Geting Instructions", "instructions", instructions)
	s.logger.Debug("END: FileService - OpenInstructionsFile")

	return instructions
}

func (s *FileService) WriteSwapFile(data []byte, dir int) {
	s.logger.Debug("INIT: FileService - WriteSwapFile")

	file, err := os.OpenFile(s.config.SwapfilePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		s.logger.Error("File Problem: Cannot Open SWAP File")

		s.logger.Debug("END: FileService - WriteSwapFile")
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteAt(data, int64(dir))
	if err != nil {
		s.logger.Error("File Problem: Cannot Write SWAP File")

		s.logger.Debug("END: FileService - WriteSwapFile")
		panic(err)
	}

	s.logger.Debug("END: FileService - WriteSwapFile")
}

func (s *FileService) ReadSwapFile(dir int) []byte {
	s.logger.Debug("INIT: FileService - ReadSwapFile")

	file, err := os.Open(s.config.SwapfilePath)
	if err != nil {
		s.logger.Error("File Problem: Cannot Open SWAP File")

		s.logger.Debug("END: FileService - ReadSwapFile")
		panic(err)
	}

	defer file.Close()

	data := make([]byte, s.config.PageSize)

	_, err = file.ReadAt(data, int64(dir))
	if err != nil {
		s.logger.Error("File Problem: Cannot Read SWAP File")

		s.logger.Debug("END: FileService - ReadSwapFile")
		panic(err)
	}

	s.logger.Debug("END: FileService - ReadSwapFile")

	return data
}

func (s *FileService) CreateDumpFile(name string, data []byte) {
	s.logger.Debug("INIT: FileService - CreateDumpFile")

	dumpFile, err := os.OpenFile(s.config.DumpPath+name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		s.logger.Error("File Problem: Cannot Create DUMP File")

		s.logger.Debug("END: FileService - CreateDumpFile")
		panic(err)
	}

	defer dumpFile.Close()

	_, err = dumpFile.WriteAt(data, 0)
	if err != nil {
		s.logger.Error("File Problem: Cannot Write DUMP File")

		s.logger.Debug("END: FileService - CreateDumpFile")
		panic(err)
	}

	s.logger.Debug("END: FileService - CreateDumpFile")
}
