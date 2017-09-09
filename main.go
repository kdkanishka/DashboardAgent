package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"strings"
	"encoding/json"
	"time"
	"./utils"
	"flag"
)

var (
	logpath = flag.String("logpath", "/tmp/dashboardagent.log", "Log Path")
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
	timeStamp := time.Now().Unix()
	if strings.Contains(frame, "InitialMonitoredService") {
		var initialMonitoredService InitialMonitoredServiceResponse
		error := json.Unmarshal([]byte(frame), &initialMonitoredService)
		if error == nil {
			fmt.Println("Successfully unmarshalled InitialMonitoredService")
		}
		for _, monitoredService := range initialMonitoredService.InitialMonitoredService.Items {
			notification := Notification{
				Name:      monitoredService.ServiceName,
				Status:    monitoredService.State,
				Item:      monitoredService,
				Timestamp: timeStamp,
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
				Name:      monitoredService.ServiceName,
				Status:    monitoredService.State,
				Item:      monitoredService,
				Timestamp: timeStamp,
			}
			notification_channel <- notification
		}
	}
}

func processNotifications(notificationsMapParam map[string]Notification, notification Notification) Command {
	prev_notification, valueExists := notificationsMapParam[notification.Name]
	if valueExists {
		prevState := prev_notification
		if notification.Status != prevState.Status {
			//state change! take necessary actions

			//check whether the service is back to its Ok state
			if notification.Status == "Ok" {
				delete(notificationsMapParam, notification.Name)
				//notify service back to normal
				return ResetAlarm{NotificationsMap: notificationsMapParam, Notification: notification}
			} else {
				notificationsMapParam[notification.Name] = notification
				//notify alarm state to the dashboard
				return PublishAlarm{NotificationsMap: notificationsMapParam, Notification: notification}
			}
		} else {
			//no state change but notification received over and over again
			notificationsMapParam[notification.Name] = notification
			return DoNothing{}
		}
	} else {
		if notification.Status != "Ok" {
			notificationsMapParam[notification.Name] = notification
			//notify alarm state to the dashboard
			return PublishAlarm{NotificationsMap: notificationsMapParam, Notification: notification}
		}
	}
	return DoNothing{}
}

func main() {
	utils.NewLog(*logpath)
	utils.Log.Println("Initializing Dashboard Agent!")
	notificationsMap := make(map[string]Notification)

	ws_channel := make(chan string)
	quite_channel := make(chan string)
	notification_channel := make(chan Notification)
	go connectToWebSocket(ws_channel, quite_channel)

	for {
		select {
		case frame := <-ws_channel:
			utils.Log.Printf("Frame receieved %s \n", frame)
			go handleFrame(frame, notification_channel)
		case notification := <-notification_channel:
			processNotifications(notificationsMap, notification)
		case sig_quite := <-quite_channel:
			utils.Log.Printf("Exit signal received! %s\n", sig_quite)
			log.Fatal("Exiting since error occured")
		}
	}
}
