package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"strings"
	"encoding/json"
)

func connectToWebSocket(ws_channel, quite_channel chan string) {
	const origin = "http://localhost:10000"
	const url = "ws://localhost:10001/"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		quite_channel <- "Unable to dial for the websocket " + err.Error()
		return
	}

	//receive frames from websocket
	for {
		var message string
		websocket.Message.Receive(ws, &message)
		if strings.Compare(message, "") == 0 {
			//empty string received this could be a result of broken pipe
			//so write some data to the websocket and verify it!
			_, err := ws.Write([]byte("test for broken pipe"))
			if err != nil {
				quite_channel <- "Connection seems to be closed. " + err.Error()
				return
			}
		} else {
			ws_channel <- message
		}
	}
}

func handleFrame(frame string, notification_channel chan Notification) {
	if strings.Contains(frame, "InitialMonitoredService") {
		var initialMonitoredService InitialMonitoredServiceResponse
		error := json.Unmarshal([]byte(frame), &initialMonitoredService)
		if error == nil {
			fmt.Println("Successfully unmarshalled InitialMonitoredService")
		}
		for _, monitoredService := range initialMonitoredService.InitialMonitoredService.Items {
			notification := Notification{
				name:   monitoredService.ServiceName,
				status: monitoredService.State,
				item:   monitoredService,
			}
			notification_channel <- notification
		}
	} else if strings.Contains(frame, "UpdatedMonitoredService") {
		var updated UpdatedMonitoredServiceResponse
		error := json.Unmarshal([]byte(frame), &updated)
		if error == nil {
			fmt.Println("Successfully unmarshalled UpdatedMonitoredService")
		}
		for _, monitoredService := range updated.UpdatedMonitoredService.Items {
			notification := Notification{
				name:   monitoredService.ServiceName,
				status: monitoredService.State,
				item:   monitoredService,
			}
			notification_channel <- notification
		}
	}
}

func main() {
	notifications := make(map[string]Notification)

	ws_channel := make(chan string)
	quite_channel := make(chan string)
	notification_channel := make(chan Notification)
	go connectToWebSocket(ws_channel, quite_channel)

	for {
		select {
		case frame := <-ws_channel:
			fmt.Printf("Frame receieved %s \n", frame)
			go handleFrame(frame, notification_channel)
		case notification := <-notification_channel:
			fmt.Println(notification)
		case sig_quite := <-quite_channel:
			fmt.Printf("Exit signal received! %s\n", sig_quite)
			log.Fatal("Exiting since error occured")
		}
	}
}
