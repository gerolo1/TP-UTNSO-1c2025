package models

import "sync"

type Process struct {
	Size                  int
	PagesAccessed         int
	InstructionsRequested int
	TimesSuspended        int
	TimesResumed          int
	TotalWrites           int
	TotalReads            int
	Mutex                 *sync.Mutex
}
