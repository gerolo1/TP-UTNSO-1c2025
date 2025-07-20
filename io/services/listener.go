package services

import (
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/io/config"
)

func StartIOListener() {
	http.HandleFunc("/io/request", HandleIORequest)

	addr := config.Config.IPIO + ":" + strconv.Itoa(config.Config.PortIO)
	config.Logger.Info("Escuchando en " + addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		config.Logger.Error("Error al iniciar servidor IO: " + err.Error())
	}
}
