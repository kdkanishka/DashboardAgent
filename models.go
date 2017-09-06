package main

type InitialMonitoredServiceResponse struct {
	InitialMonitoredService InitialMonitoredService `json:"InitialMonitoredService"`
}

type InitialMonitoredService struct {
	Items []InitialMonitoredServiceItem `json:"items"`
}

type InitialMonitoredServiceItem struct {
	ServiceName string `json:"serviceName"`
	State string `json:"state"`
	Alarms []Alarm `json:"alarms"`
}

type Alarm struct {
	ResponsibleUnit string `json:"responsibleUnit"`
	Label string `json:"label"`
	AlarmType string `json:"alarmType"`
	Origin string `json:"origin"`
	Details string `json:"details"`
	Severity string `json:severity`
}


