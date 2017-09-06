package main

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestMarshalAlarm(t *testing.T) {
	alarm := Alarm{
		ResponsibleUnit: "Kingfisher",
		Label:           "Cannot Keep up!",
		AlarmType:       "QueueAlarm",
		Origin:          "archiveindex-service",
		Details:         "Consumer can't keep up",
		Severity:        "Warning",
	}

	data, error := json.Marshal(alarm)
	if error != nil {
		t.Error("Unable to marshal data!")
	}
	t.Log(string(data))
}

func TestUnmarshalAlarm(t *testing.T) {
	jsonText := `{"responsibleUnit":"Kingfisher","label":"Cannot Keep up!","alarmType":"QueueAlarm","origin":"archiveindex-service","details":"Consumer can't keep up","Severity":"Warning"}`
	var alarm Alarm
	error := json.Unmarshal([]byte(jsonText), &alarm)
	if error != nil {
		t.Error("Unable to unmarshal json encoded text to Alarm structure")
	}

	responsibleUnit := "Kingfisher"
	if alarm.ResponsibleUnit != responsibleUnit {
		missingMember(t, "ResponsibleUnit", responsibleUnit)
	}

	label := "Cannot Keep up!"
	if alarm.Label != label {
		missingMember(t, "Label", label)
	}

	alarmType := "QueueAlarm"
	if alarm.AlarmType != alarmType {
		missingMember(t, "AlarmType", alarmType)
	}

	origin := "archiveindex-service"
	if alarm.Origin != origin {
		missingMember(t, "Origin", origin)
	}

	details := "Consumer can't keep up"
	if alarm.Details != details {
		missingMember(t,"Details",details)
	}

	severity := "Warning"
	if alarm.Severity != severity {
		missingMember(t,"Severity", severity)
	}
}

func missingMember(t *testing.T, member string, value string) {
	t.Error(fmt.Sprintf("Member '%s' doesn't contain proper value. Value found was : '%s'",
		member, value))
}
