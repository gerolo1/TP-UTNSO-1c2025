package config

import (
	"encoding/json"
	"os"

	"github.com/sisoputnfrba/tp-golang/io/models"
)

var Config *models.Config

func InitConfiguration() {
	configFile, e := os.Open("io/config/config.json")
	if e != nil {
		panic(e)
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	e = jsonParser.Decode(&Config)
	if e != nil {
		panic(e)
	}
}
