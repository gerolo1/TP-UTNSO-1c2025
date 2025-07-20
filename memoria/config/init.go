package config

import (
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"sync"
)

func InitUserMemory(memorySize int) []byte {
	return make([]byte, memorySize)
}

func InitFreeMemory(memorySize int) *int {
	freeMemory := new(int)
	freeMemory = &memorySize
	return freeMemory
}

func InitFrameBitMap(memorySize int, pageSize int) []bool {
	return make([]bool, memorySize/pageSize)
}

func InitMutexUserMemory() *sync.Mutex {
	return &sync.Mutex{}
}

func InitMutexFreeMemory() *sync.Mutex {
	return &sync.Mutex{}
}

func InitMutexFrameBitMap() *sync.Mutex {
	return &sync.Mutex{}
}

////

func InitSwapBitMap() []bool {
	return []bool{}
}

func InitSwapEntries() *[]*models.SwapEntry {
	return &[]*models.SwapEntry{}
}

func InitMutexSwapBitMap() *sync.Mutex {
	return &sync.Mutex{}
}

func InitMutexSwapEntries() *sync.Mutex {
	return &sync.Mutex{}
}

/////

func InitProcesses() map[string]*models.Process {
	return make(map[string]*models.Process)
}

func InitMutexProcesses() *sync.RWMutex {
	return &sync.RWMutex{}
}

func InitProcessesInstructions() map[string]*models.ProcessInstructions {
	return make(map[string]*models.ProcessInstructions)
}

func InitMutexProcessesInstructions() *sync.RWMutex {
	return &sync.RWMutex{}
}

func InitProcessesPages() map[string]*models.ProcessPages {
	return make(map[string]*models.ProcessPages)
}

func InitMutexProcessesPages() *sync.RWMutex {
	return &sync.RWMutex{}
}
