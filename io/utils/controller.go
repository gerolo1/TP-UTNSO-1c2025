package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/io/config"
)

func SendToModule(targetIP string, targetPort int, endpoint string, payload interface{}, v interface{}) error {
	// Convert payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		config.Logger.Error("Failed to encode payload")
		return err
	}

	// Build the target URL
	url := fmt.Sprintf("http://%s:%d/%s", targetIP, targetPort, endpoint)

	// Send the POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		config.Logger.Error("Failed to send request to" + url)
		return err
	}

	//Cierra la peitición cuando termina la función
	defer resp.Body.Close()

	config.Logger.Debug("POST " + url + " status " + resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil
	} else {
		if v != nil {
			return json.NewDecoder(resp.Body).Decode(v)
		} else {
			return nil
		}
	}
}

func GetFromModule(targetIP string, targetPort int, endpoint string, v interface{}) error {
	url := fmt.Sprintf("http://%s:%d/%s", targetIP, targetPort, endpoint)

	resp, err := http.Get(url)
	if err != nil {
		config.Logger.Error("Error al hacer la petición GET " + url)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		config.Logger.Error("el módulo respondió con código " + strconv.Itoa(resp.StatusCode))
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("falló al leer el cuerpo de la respuesta")
		return err
	}

	config.Logger.Debug("GET " + url + " status " + resp.Status)

	return json.Unmarshal(body, v)
}
