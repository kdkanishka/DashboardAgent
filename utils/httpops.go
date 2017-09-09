package utils

import (
	"net/http"
	"bytes"
	"io/ioutil"
)

func Post(data []byte, url string) {
	response, error := http.Post(url, "application/json", bytes.NewBuffer(data))
	defer response.Body.Close()

	if error != nil || response.StatusCode != 200 {
		Log.Printf("Unable to POST data to %s , %s\n", url, error.Error())
		Log.Printf("Status code %d \n", response.StatusCode)
	}

	respo, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		Log.Printf("Unable to read from the HTTP response!, %s \n", readErr.Error())
	}
	Log.Printf("Response from service %s\n", string(respo))
}
