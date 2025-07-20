package models

type TLBEntry struct {
    Page        int
    Frame         int
    LastAccess  int
    EntryOrder  int
}

type TLB struct {
    Entries      []TLBEntry
    MaxEntries   int
    OrderCont int // para FIFO
    AccesCont int // para LRU
}
