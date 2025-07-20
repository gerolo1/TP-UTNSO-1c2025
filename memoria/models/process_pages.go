package models

import "sync"

type ProcessPages struct {
	Mutex            *sync.Mutex
	Pages            *[]*Page
	ContadorEntradas int
}
