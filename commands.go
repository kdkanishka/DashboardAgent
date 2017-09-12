package main

import (
	"encoding/json"
	"./utils"
)

func kradiatorNotificationEndpoint() string {
	config := utils.Read()
	return config.KradiatorEndpoint
}

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
	utils.Post(serialized, kradiatorNotificationEndpoint() + "/PublishAlarm")
}

func (command ResetAlarm) exec() {
	utils.Log.Println("Executing ResetAlarm")
	//encode in json
	serialized, error := json.Marshal(command)
	if error != nil {
		utils.Log.Println("Unable to marshal ResetAlarm object," + error.Error())
	}

	//post
	utils.Post(serialized, kradiatorNotificationEndpoint() + "/ResetAlarm")
}

func (Command DoNothing) exec() {
}
