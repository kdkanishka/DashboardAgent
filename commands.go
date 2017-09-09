package main

import (
	"encoding/json"
	"./utils"
)

var kradiatorNotificationEndpoint = "http://10.2.4.33:8000"

type Command interface {
	exec()
}

type PublishAlarm struct {
	NotificationsMap map[string]Notification
	Notification     Notification
}

type ResetAlarm struct {
	NotificationsMap map[string]Notification
	Notification     Notification
}

type DoNothing struct {
}

func (command PublishAlarm) exec() {
	//encode in json
	serialized, error := json.Marshal(command)
	if error != nil {
		utils.Log.Println("Unable to marshal PublishAlarm object," + error.Error())
	}

	//post
	utils.Post(serialized, kradiatorNotificationEndpoint)
}

func (command ResetAlarm) exec() {
	//encode in json
	serialized, error := json.Marshal(command)
	if error != nil {
		utils.Log.Println("Unable to marshal ResetAlarm object," + error.Error())
	}

	//post
	utils.Post(serialized, kradiatorNotificationEndpoint)
}

func (Command DoNothing) exec() {
}
