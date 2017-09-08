package main

import (
	"testing"
	"fmt"
)

func TestAlarmShouldBePublishedOnCriticalSeverity(t *testing.T) {
	notificationsMap := make(map[string]Notification)
	notificationsMap["service1"] = Notification{name: "service1", status: "Error"}
	notificationsMap["service2"] = Notification{name: "service2", status: "Error"}
	notificationsMap["service3"] = Notification{name: "service3", status: "Error"}

	newNotification := Notification{name: "service4",status: "Error"}

	resultingCommand := processNotifications(notificationsMap,newNotification)
	val, isCorrectType := resultingCommand.(PublishAlarm)
	fmt.Println(val)
	if !isCorrectType {
		t.Error("Incorrect command returned")
	}

	_, isExist := notificationsMap["service4"]
	if !isExist {
		t.Log(notificationsMap)
		t.Error("New service with Error status could not find in notifications map")
	}
}