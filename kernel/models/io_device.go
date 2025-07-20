package models

type IoDevice struct {
	Queue     []*IoQueueEntry
	Instances []*IoInstance
}
