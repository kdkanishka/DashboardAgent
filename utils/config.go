package utils

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	KradiatorEndpoint string `json:"kradiatorEndpoint"`
}

func Read() Configuration {
	fileContent, error := ioutil.ReadFile("config.json")
	if error != nil {
		Log.Printf("Unable to open configuration file %s", error.Error())
		panic("Unable to continue! since configuration file is not available")
	}

	configuration := Configuration{}
	decodeError := json.Unmarshal(fileContent, &configuration)

	if decodeError != nil {
		Log.Printf("Unable to decode config file content %s", decodeError.Error())
		panic("Unable to continue! since configuration problem")
	}
	return configuration
}
