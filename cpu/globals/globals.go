package globals

import (
	"github.com/sisoputnfrba/tp-golang/cpu/models"
	"github.com/sisoputnfrba/tp-golang/cpu/models/dtos"
	"sync"
)

var (
	ConfigFile      *models.Config
	PageSize        int
	Tlb             *models.TLB
	Cache           *models.Cache
	Identifier      string
	PC              int
	PID             int
	Levels          int
	EntriesPerTable int
)

var InterruptMutex sync.Mutex
var InterruptSync *sync.Mutex
var InterruptExist bool = false
var Interrupt dtos.Interrupt
