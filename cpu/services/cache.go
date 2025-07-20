package services

import (
	"strconv"
	"strings"

	"github.com/sisoputnfrba/tp-golang/cpu/config"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/models"
	"github.com/sisoputnfrba/tp-golang/cpu/models/dtos"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
)

// Busca una página en la caché.
// Si la encuentra, devuelve su contenido y actualiza el bit de uso.
func GetFromCache(page int) (string, bool, int) {
	for i := range globals.Cache.Entries {
		entry := &globals.Cache.Entries[i]
		if entry.Page == page {
			entry.UseBit = true
			config.Logger.Info("CACHE HIT - Página " + strconv.Itoa(page))
			return entry.Content, true, entry.Frame
		}
	}
	//config.Logger.Info("CACHE MISS - Página " + strconv.Itoa(page))
	return "", false, 0
}

func writeCache(page int, content string) {
	for i := range globals.Cache.Entries {
		entry := &globals.Cache.Entries[i]
		if entry.Page == page {
			entry.Content = content
			entry.UseBit = true
			entry.ModifiedBit = true
		}
	}
}

func readCache(pageContent string, offset int, size int) string {
	end := offset + size
	if end > len(pageContent) {
		end = len(pageContent)
	}
	result := pageContent[offset:end]
	return result
}

// Agrega una página a la caché utilizando el algoritmo CLOCK
// Si se reemplaza una entrada modificada, la escribe de vuelta en memoria.
func AddToCache(page int, content string, modified bool, frame int) {
	cache := globals.Cache

	// Si hay espacio se agrega directo
	if len(cache.Entries) < cache.MaxEntries {
		cache.Entries = append(cache.Entries, models.CacheEntry{
			Page:        page,
			Frame:       frame, // asigno frame acá
			Content:     content,
			UseBit:      true,
			ModifiedBit: modified,
		})
		config.Logger.Debug("Página agregada a la caché: " + strconv.Itoa(page))
		return
	}

	replacement := strings.ToUpper(globals.ConfigFile.CacheReplacement)

	// Algoritmos de reemplazo cuando no hay espacio

	// CLOCK
	if replacement == "CLOCK" {
		for {
			ptr := cache.ClockPointer
			entry := &cache.Entries[ptr]

			if !entry.UseBit {
				if entry.ModifiedBit {
					err := WritePageBackToMemory(entry.Page, entry.Content, entry.Frame)
					if err != nil {
						config.Logger.Error("Error escribiendo página modificada a memoria: " + err.Error())
					}
				}

				*entry = models.CacheEntry{
					Page:        page,
					Frame:       frame, // asigno frame acá
					Content:     content,
					UseBit:      true,
					ModifiedBit: modified,
				}

				cache.ClockPointer = (ptr + 1) % cache.MaxEntries
				config.Logger.Debug("Página reemplazada en caché (CLOCK): " + strconv.Itoa(page))
				return
			}

			entry.UseBit = false
			cache.ClockPointer = (ptr + 1) % cache.MaxEntries
		}
	}

	// CLOCK-M (modificado)
	if replacement == "CLOCK-M" {
		// Primera pasada: buscar UseBit=false y ModifiedBit=false
		for i := 0; i < cache.MaxEntries; i++ {
			ptr := cache.ClockPointer
			entry := &cache.Entries[ptr]

			if !entry.UseBit && !entry.ModifiedBit {
				*entry = models.CacheEntry{
					Page:        page,
					Frame:       frame, // asigno frame acá
					Content:     content,
					UseBit:      true,
					ModifiedBit: modified,
				}
				cache.ClockPointer = (ptr + 1) % cache.MaxEntries
				config.Logger.Debug("Página reemplazada en caché (CLOCK-M - limpio): " + strconv.Itoa(page))
				return
			}

			cache.ClockPointer = (ptr + 1) % cache.MaxEntries
		}

		// Segunda pasada: buscar UseBit=false y ModifiedBit=true
		for i := 0; i < cache.MaxEntries; i++ {
			ptr := cache.ClockPointer
			entry := &cache.Entries[ptr]

			if !entry.UseBit && entry.ModifiedBit {
				err := WritePageBackToMemory(entry.Page, entry.Content, entry.Frame)
				if err != nil {
					config.Logger.Error("Error escribiendo página modificada a memoria (CLOCK-M): " + err.Error())
				}

				*entry = models.CacheEntry{
					Page:        page,
					Frame:       frame, // <-- asigno frame acá
					Content:     content,
					UseBit:      true,
					ModifiedBit: modified,
				}
				cache.ClockPointer = (ptr + 1) % cache.MaxEntries
				config.Logger.Debug("Página reemplazada en caché (CLOCK-M - modificada): " + strconv.Itoa(page))
				return
			}

			cache.ClockPointer = (ptr + 1) % cache.MaxEntries
		}

		// Tercera pasada: limpiar todos los UseBit y volver a intentar
		for i := range cache.Entries {
			cache.Entries[i].UseBit = false
		}

		// Volver a primera pasada
		AddToCache(page, content, modified, frame)
	}
}

// Envía a memoria la página que fue modificada antes de ser reemplazada
func WritePageBackToMemory(page int, content string, frame int) error {
	df := frame // Dirección física base del marco

	dto := dtos.WriterDTO{
		PID:   globals.PID,
		DF:    df,
		Value: []byte(content),
	}

	err := utils.SendToModule(
		globals.ConfigFile.IPMemory,
		globals.ConfigFile.PortMemory,
		"memory/write",
		dto,
		nil,
	)
	if err != nil {
		config.Logger.Error("Error al escribir en memoria cache actualizada. ")
	} else {
		//“PID: <PID> - Memory Update - Página: <NUMERO_PAGINA> - Frame: <FRAME_EN_MEMORIA_PRINCIPAL>
		config.Logger.Info("PID: " + strconv.Itoa(globals.PID) + " - Memory Update - Página: " + strconv.Itoa(page) + " - Frame: " + strconv.Itoa(frame))
	}
	return err
}

func FlushModifiedPagesToMemory(pid int) {
	for _, entry := range globals.Cache.Entries {
		if entry.ModifiedBit {
			err := WritePageBackToMemory(entry.Page, entry.Content, entry.Frame)
			if err != nil {
				config.Logger.Error("Error al escribir página modificada a memoria (flush): " + err.Error())
			} else {
				config.Logger.Debug("Página modificada enviada a memoria en flush: " + strconv.Itoa(entry.Page))
				config.Logger.Info("PID: " + strconv.Itoa(pid) + " - Memory update - Página:" + strconv.Itoa(entry.Page) + " - Frame: " + strconv.Itoa(entry.Frame))
			}
		}
	}
}
