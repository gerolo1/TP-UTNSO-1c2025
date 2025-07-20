package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendToModule(ip string, port int, endpoint string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error al serializar payload: %v", err)
		return err
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error al enviar request a %s: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("POST %s - Status: %s", url, resp.Status)
	return nil
}
