package models

type CacheEntry struct {
	Page          int
	Frame         int
	Content       string
	UseBit        bool
	ModifiedBit   bool
}

type Cache struct {
	Entries      []CacheEntry
	MaxEntries   int
	ClockPointer  int
}