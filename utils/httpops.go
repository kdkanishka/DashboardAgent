package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(apikey string, url string) {
	heartBeatURL := fmt.Sprintf("%s/update?api_key=%s&field1=%d", url, apikey, 1)
	response, error := http.Get(heartBeatURL)
	defer func() {
		if error == nil {
			response.Body.Close()
		}
	}()

	if error != nil || response.StatusCode != 200 {
		Log.Printf("Unable to publish heartbeart! %s\n", heartBeatURL)
		Log.Printf("Error %s \n", error.Error())
	} else {
		respo, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			Log.Printf("Unable to read from the HTTP response!, %s \n", readErr.Error())
		}
		Log.Printf("Response from service %s\n", string(respo))
	}
}

func Post(data []byte, url string, verbose bool) {
	response, error := http.Post(url, "application/json", bytes.NewBuffer(data))
	defer func() {
		if error == nil {
			response.Body.Close()
			if verbose {
				Log.Printf("Response closed.")
			}
		}
	}()

	if error != nil || response.StatusCode != 200 {
		Log.Printf("Unable to POST data to %s\n", url)
		Log.Printf("Error %s \n", error.Error())
	} else {
		respo, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			Log.Printf("Unable to read from the HTTP response!, %s \n", readErr.Error())
		}
		if verbose {
			Log.Printf("Response from service %s\n", string(respo))
		}
	}
}
