package main

import (
	"testing"
)

func Test_alarm_should_be_published_on_critical_severity_for_new_failure(t *testing.T) {
	notificationsMap := make(map[string]Notification)
	notificationsMap["service1"] = Notification{Name: "service1", Status: "Error", Timestamp: 1}
	notificationsMap["service2"] = Notification{Name: "service2", Status: "Error", Timestamp: 1}
	notificationsMap["service3"] = Notification{Name: "service3", Status: "Error", Timestamp: 1}

	newNotification := Notification{Name: "service4", Status: "Error", Timestamp: 2}

	resultingCommand := processNotifications(notificationsMap, newNotification)
	_, isCorrectType := resultingCommand.(PublishAlarm)
	if !isCorrectType {
		t.Error("Incorrect command returned")
	}

	notification, isExist := notificationsMap["service4"]
	if !isExist {
		t.Error("New service with Error Status could not find in notifications map")
	}

	if notification.Timestamp != 2 {
		t.Error("New notification hasn't been updated properly")
	}
}

func Test_Alarm_should_not_be_published_on_critical_severity_for_known_failure_but_notficiationmap_should_be_updated(t *testing.T) {
	notificationsMap := make(map[string]Notification)
	notificationsMap["service1"] = Notification{Name: "service1", Status: "Error", Timestamp: 1}
	notificationsMap["service2"] = Notification{Name: "service2", Status: "Error", Timestamp: 1}
	notificationsMap["service3"] = Notification{Name: "service3", Status: "Error", Timestamp: 1}

	newNotification := Notification{Name: "service3", Status: "Error", Timestamp: 2}

	resultingCommand := processNotifications(notificationsMap, newNotification)
	_, isCorrectType := resultingCommand.(DoNothing)
	if !isCorrectType {
		t.Error("Incorrect command returned")
	}

	notification, isExist := notificationsMap["service3"]
	if ! isExist {
		t.Error("There should be a known notification for the service")
	}

	if notification.Timestamp != 2 {
		t.Error("New notification hasn't been updated properly")
	}
}

func Test_Alarm_should_be_published_on_critical_severity_for_known_failure_with_none_Ok_status(t *testing.T){
	notificationsMap := make(map[string]Notification)
	notificationsMap["service1"] = Notification{Name: "service1", Status: "Error", Timestamp: 1}
	notificationsMap["service2"] = Notification{Name: "service2", Status: "Error", Timestamp: 1}
	notificationsMap["service3"] = Notification{Name: "service3", Status: "Error", Timestamp: 1}

	newNotification := Notification{Name: "service3", Status: "Critical", Timestamp: 2}

	resultingCommand := processNotifications(notificationsMap, newNotification)
	_, isCorrectType := resultingCommand.(PublishAlarm)
	if !isCorrectType {
		t.Error("Incorrect command returned")
	}

	notification, isExist := notificationsMap["service3"]
	if ! isExist {
		t.Error("There should be a known notification for the service")
	}

	if notification.Timestamp != 2 {
		t.Error("New notification hasn't been updated properly")
	}

	if notification.Status != "Critical" {
		t.Error("Notification Status hasn't been updated properly")
	}
}

func Test_Ok_Notification_for_known_service_should_result_in_ResetAlarm_and_notificationmap_should_be_updated(t *testing.T) {
	notificationsMap := make(map[string]Notification)
	notificationsMap["service1"] = Notification{Name: "service1", Status: "Error", Timestamp: 1}
	notificationsMap["service2"] = Notification{Name: "service2", Status: "Error", Timestamp: 1}
	notificationsMap["service3"] = Notification{Name: "service3", Status: "Error", Timestamp: 1}

	newNotification := Notification{Name: "service3", Status: "Ok", Timestamp: 2}

	resultingCommand := processNotifications(notificationsMap, newNotification)
	_, isCorrectType := resultingCommand.(ResetAlarm)
	if !isCorrectType {
		t.Error("Incorrect command returned")
	}

	_, isExist := notificationsMap["service3"]
	if isExist {
		t.Error("Service should be removed from the notification map after receiving Ok Status")
	}
}

func Test_Ok_Notification_for_unknown_service_should_result_in_DoNothing(t *testing.T) {
	notificationsMap := make(map[string]Notification)
	notificationsMap["service1"] = Notification{Name: "service1", Status: "Error", Timestamp: 1}
	notificationsMap["service2"] = Notification{Name: "service2", Status: "Error", Timestamp: 1}
	notificationsMap["service3"] = Notification{Name: "service3", Status: "Error", Timestamp: 1}

	newNotification := Notification{Name: "service4", Status: "Ok", Timestamp: 2}

	resultingCommand := processNotifications(notificationsMap, newNotification)
	_, isCorrectType := resultingCommand.(DoNothing)
	if !isCorrectType {
		t.Error("Incorrect command returned")
	}

	_, isExist := notificationsMap["service4"]
	if isExist {
		t.Error("No service entry should be added to the notification map if it is unknown")
	}
}
