package utils

import (
	"encoding/json"
	"io/ioutil"
	"flag"
)

type Configuration struct {
	KradiatorEndpoint string `json:"kradiatorEndpoint"`
}

func Read() Configuration {
	flagFile := flag.String("config", "config.json", "Configuration file")
	flag.Parse()

	fileContent, error := ioutil.ReadFile(*flagFile)
	if error != nil {
		Log.Printf("Unable to open configuration file %s File:%s\n", error.Error(), *flagFile)
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
