package models

import "time"

type CpuModule struct {
	Id           string
	Address      string
	EstNextBurst float64
	InitDate     time.Time
	Pcb          *Pcb
	Busy         bool
	Interrupt    bool
}
