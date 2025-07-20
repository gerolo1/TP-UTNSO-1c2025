package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

func InitServer(w http.ResponseWriter, r *http.Request) {
	log.Println("kernel server ok")

	// el cliente sabe que recibira JSON
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode("OK")
	if err != nil {
		log.Println("Error codificando JSON", err)
		return
	}
}
