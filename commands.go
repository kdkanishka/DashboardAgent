package main

import (
	"encoding/json"

	"github.com/kdkanishka/DashboardAgent/utils"
)

type Command interface {
	exec()
}

type PublishAlarm struct {
	NotificationsMap map[string]Notification `json:"notificationsMap"`
	Notification     Notification            `json:"notification"`
}

type ResetAlarm struct {
	NotificationsMap map[string]Notification `json:"notificationsMap"`
	Notification     Notification            `json:"notification"`
}

type PublishClockNotification struct {
	ClockValue string `json:"clockValue"`
}

type DoNothing struct {
}

func (command PublishAlarm) exec() {
	utils.Log.Println("Executing PublishAlarm")
	//encode in json
	serialized, error := json.Marshal(command)
	if error != nil {
		utils.Log.Println("Unable to marshal PublishAlarm object," + error.Error())
	}

	//post
	utils.Post(serialized, utils.KradiatorNotificationEndpoint()+"/PublishAlarm", true)
}

func (command ResetAlarm) exec() {
	utils.Log.Println("Executing ResetAlarm")
	//encode in json
	serialized, error := json.Marshal(command)
	if error != nil {
		utils.Log.Println("Unable to marshal ResetAlarm object," + error.Error())
	}

	//post
	utils.Post(serialized, utils.KradiatorNotificationEndpoint()+"/ResetAlarm", true)
}

func (command PublishClockNotification) exec() {
	//utils.Log.Println("Executing PublishClockNotification")
	//encode in json
	serialized, error := json.Marshal(command)
	if error != nil {
		utils.Log.Println("Unable to marshal PublishClockNotification object," + error.Error())
	}

	//post
	utils.Post(serialized, utils.KradiatorNotificationEndpoint()+"/PublishClockNotification", false)
}

func (Command DoNothing) exec() {
}
