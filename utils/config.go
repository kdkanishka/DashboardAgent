package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configuration struct {
	KradiatorEndpoint string `json:"kradiatorEndpoint"`
}

func KradiatorNotificationEndpoint() string {
	config := Read()
	return config.KradiatorEndpoint
}

func Read() Configuration {
	configFile := os.Args[1]

	fileContent, error := ioutil.ReadFile(configFile)
	if error != nil {
		Log.Printf("Unable to open configuration file %s File:%s\n", error.Error(), configFile)
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
