package main

import (
	"encoding/json"
	"flag"
	"log"
	"strings"
	"time"
	"net"

	"github.com/kdkanishka/DashboardAgent/utils"
	"golang.org/x/net/websocket"
)

var (
	logpath = flag.String("logpath", "/tmp/dashboardagent.log", "Log Path")
)

func connectToWebSocket(wsChannel, quiteChannel chan string) {
	const origin = "http://localhost:10000"
	const url = "ws://localhost:10000/ws"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		quiteChannel <- "Unable to dial for the websocket " + err.Error()
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
				quiteChannel <- "Connection seems to be closed. " + err.Error()
				return
			}
		} else {
			wsChannel <- message
		}
	}
}

func handleFrame(frame string, notificationChannel chan Notification) {
	timeStamp := time.Now().Unix()
	if strings.Contains(frame, "InitialMonitoredService") {
		var initialMonitoredService InitialMonitoredServiceResponse
		error := json.Unmarshal([]byte(frame), &initialMonitoredService)
		if error == nil {
			utils.Log.Println("Successfully unmarshalled InitialMonitoredService")
		}
		for _, monitoredService := range initialMonitoredService.InitialMonitoredService.Items {
			notification := Notification{
				Name:      monitoredService.ServiceName,
				Status:    monitoredService.State,
				Item:      monitoredService,
				Timestamp: timeStamp,
			}
			notificationChannel <- notification
		}
	} else if strings.Contains(frame, "UpdatedMonitoredService") {
		var updated UpdatedMonitoredServiceResponse
		error := json.Unmarshal([]byte(frame), &updated)
		if error == nil {
			utils.Log.Println("Successfully unmarshalled UpdatedMonitoredService")
		}
		for _, monitoredService := range updated.UpdatedMonitoredService.Items {
			notification := Notification{
				Name:      monitoredService.ServiceName,
				Status:    monitoredService.State,
				Item:      monitoredService,
				Timestamp: timeStamp,
			}
			notificationChannel <- notification
		}
	} else if strings.Contains(frame, "Clock") {
		var clockFrame ClockFrame
		error := json.Unmarshal([]byte(frame), &clockFrame)
		if error == nil {
			//utils.Log.Println("Successfully unmarshalled Clock") //TODO remove this line
			//utils.Log.Println(clockFrame.Clock.Time)
			notification := Notification{
				IsClockNotification: true,
				Name:                "Clock",
				Status:              clockFrame.Clock.Time,
			}
			notificationChannel <- notification
		}
	}
}

func processNotifications(notificationsMapParam map[string]Notification, notification Notification) Command {
	//check whether that this is just a clock notification (clock updated receives 6 times per minute)
	if notification.IsClockNotification == true {
		return PublishClockNotification{ClockValue: notification.Status}
	} else {
		//actual dashbard notification
		prevNotification, valueExists := notificationsMapParam[notification.Name]
		if valueExists {
			prevState := prevNotification
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
	}
	return DoNothing{}
}

func main() {
	utils.NewLog(*logpath)
	utils.Log.Println("Initializing Dashboard Agent!")
	notificationsMap := make(map[string]Notification)

	wsChannel := make(chan string)
	quiteChannel := make(chan string)
	notificationChannel := make(chan Notification)

	go connectToWebSocket(wsChannel, quiteChannel)
	//go heartBeatScheduler()

	for {
		select {
		case frame := <-wsChannel:
			//utils.Log.Printf("Frame receieved %s \n", frame)
			go handleFrame(frame, notificationChannel)
		case notification := <-notificationChannel:
			resultingCommand := processNotifications(notificationsMap, notification)
			resultingCommand.exec()
		case sigQuite := <-quiteChannel:
			utils.Log.Printf("Exit signal received! %s\n", sigQuite)
			log.Fatal("Exiting since error occured")
		}
	}
}
