package models

import "sync"

type ProcessInstructions struct {
	Mutex        *sync.Mutex
	Instructions []string
}
