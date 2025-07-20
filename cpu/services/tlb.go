package services

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/models"
)

func AddEntry(page int, frame int, replacementPolicy string) {
	tlb := globals.Tlb
	if len(tlb.Entries) >= tlb.MaxEntries {
		switch replacementPolicy {
		case "FIFO":
			tlb.Entries = tlb.Entries[1:]
		case "LRU":
			lruIdx := 0
			minAccess := tlb.Entries[0].LastAccess
			for i, entry := range tlb.Entries {
				if entry.LastAccess < minAccess {
					minAccess = entry.LastAccess
					lruIdx = i
				}
			}
			// Elimina el LRU
			tlb.Entries = append(tlb.Entries[:lruIdx], tlb.Entries[lruIdx+1:]...)
		}
	}

	// Agregar la nueva entrada
	tlb.Entries = append(tlb.Entries, models.TLBEntry{
		Page:       page,
		Frame:      frame,
		LastAccess: tlb.AccesCont,
		EntryOrder: tlb.OrderCont,
	})

	tlb.OrderCont++
	tlb.AccesCont++
}

func GetFrameFromLocalTLB(page int) (int, bool) {
	tlb := globals.Tlb
	for i := 0; i < len(tlb.Entries); i++ {
		entry := &tlb.Entries[i]
		if entry.Page == page {
			entry.LastAccess = tlb.AccesCont
			tlb.AccesCont++
			return entry.Frame, true
		}
	}
	return -1, false
}
