package models

type Pcb struct {
	Pid          int
	Pc           int
	Size         int
	File         string
	StateMap     map[string]*State
	CurrentState string
	EstPrevBurst float64
	PrevBurst    float64
	NextBurst    float64
}
