package models

type Page struct {
	ID      int
	Level   int
	Entries []*Entry
}
