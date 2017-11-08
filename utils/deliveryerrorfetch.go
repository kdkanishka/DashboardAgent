package utils

import (
	"time"
	"errors"
	"os/exec"
	"bufio"
	"strings"
	"strconv"
	"encoding/json"
	"fmt"
)

type DeliveryServiceSendErrors struct {
	Data map[string]map[string]int `json:"data"`
}

func FetchDeliveryErrors() {
	Log.Printf("Fetching delivery errors..")

	command1 := "ssh thn-prod-service4 'cd /var/log/services/delivery1/old_logs/;for X in *; do printf $X; printf \" \"; zgrep \"SEND_ERROR\" $X | wc -l; done'"
	command2 := "ssh thn-prod-service7 'grep \"SEND_ERROR\" /var/log/services/delivery1/delivery-service.log | wc -l'"

	command3 := "ssh thn-prod-service4 'cd /var/log/services/delivery2/old_logs/;for X in *; do printf $X; printf \" \"; zgrep \"SEND_ERROR\" $X | wc -l; done'"
	command4 := "ssh thn-prod-service7 'grep \"SEND_ERROR\" /var/log/services/delivery2/delivery-service.log | wc -l'"

	deliveryErrMap1, err1 := getErrorLogsForService(command1, command2)
	deliveryErrMap2, err2 := getErrorLogsForService(command3, command4)

	if err1 == nil && err2 == nil {
		Log.Println("Successfully completed fetching delivery errors!")
		deliveryServiceSendErros := DeliveryServiceSendErrors{Data: make(map[string]map[string]int)}
		deliveryServiceSendErros.Data["delivery1"] = deliveryErrMap1
		deliveryServiceSendErros.Data["delivery2"] = deliveryErrMap2

		serialized, jsonErr := json.Marshal(deliveryServiceSendErros)
		if jsonErr == nil {
			Log.Println(string(serialized))
			//POST data to KRadiator
			Post(serialized, KradiatorNotificationEndpoint()+"/DeliveryServiceSendErrors")
		} else {
			Log.Println(jsonErr)
		}
	} else {
		Log.Println("Failed to fetch delivery errors!")
		if err1 != nil {
			Log.Println(err1.Error())
		}
		if err2 != nil {
			Log.Println(err2.Error())
		}
	}
}

func getErrorLogsForService(commands ...string) (map[string]int, error) {
	commandOutputChan := make(chan string)
	errFlagChan := make(chan bool)

	for _, command := range commands {
		go execProcess(command, commandOutputChan, errFlagChan)
	}

	return compileFinalOtput(commandOutputChan, errFlagChan, len(commands))
}

func compileFinalOtput(output chan string, errFlagChan chan bool, expectedOutputCOunt int) (map[string]int, error) {
	Log.Println("Preparing final output..")
	dataMap := make(map[string]int)
	for i := 0; i < expectedOutputCOunt; i++ {
		select {
		case output := <-output:
			Log.Println("Processing today's error count")
			if parsedCount, err := strconv.Atoi(strings.TrimSpace(output)); err == nil {
				//recieved todays error count
				year, month, day := time.Now().Date()
				today := fmt.Sprintf("%d-%d-%d", year, month, day)
				dataMap[today] = parsedCount
			} else {
				Log.Println("Processing old logs data")
				fetchedMap := logFileDataToMap(output)
				for k, v := range fetchedMap {
					dataMap[k] = v
				}
			}
		case <-errFlagChan:
			Log.Println("unable to execute command successfully")
			return nil, errors.New("unable to execute command successfully")
		case <-time.After(time.Minute * 3):
			Log.Println("unable to wait until commands complete, timeout occured")
			return nil, errors.New("unable to wait until commands complete, timeout occured")
		}
	}
	Log.Println("Prepared final output successfully!")
	return dataMap, nil
}

func execProcess(command string, output chan string, errFlag chan bool) {
	Log.Printf("Executing command : %s\n", command)
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.Output()
	//stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		Log.Printf("Error occured while executing command %s \n Error : %s", command, err.Error())
		errFlag <- true
	}
	commandOutput := string(out)
	Log.Printf("Done with command %s \n Command output %s", command, commandOutput)
	output <- commandOutput
}

func logFileDataToMap(logData string) map[string]int {
	Log.Println("Converting old log file data to a map")
	dataMap := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(logData))
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		key := strings.Replace(strings.Replace(line, "delivery-service.log-", "", 1), ".zip", "", 1)
		splitted := strings.Split(key, " ")

		errorCount, err := strconv.Atoi(splitted[1])
		if err == nil {
			dataMap[splitted[0]] = errorCount
		}
	}
	Log.Println("Converted old log file data to a map successfully!")
	return dataMap
}
