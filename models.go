package main

type InitialMonitoredServiceResponse struct {
	InitialMonitoredService InitialMonitoredService `json:"InitialMonitoredService"`
}

type UpdatedMonitoredServiceResponse struct {
	UpdatedMonitoredService UpdatedMonitoredService `json:"UpdatedMonitoredService"`
}

type UpdatedMonitoredService struct {
	Items []InitialMonitoredServiceItem `json:"items"`
}

type InitialMonitoredService struct {
	Items []InitialMonitoredServiceItem `json:"items"`
}

type InitialMonitoredServiceItem struct {
	ServiceName string  `json:"serviceName"`
	State       string  `json:"state"`
	Alarms      []Alarm `json:"alarms"`
}

type Alarm struct {
	ResponsibleUnit string `json:"responsibleUnit"`
	Label           string `json:"label"`
	AlarmType       string `json:"alarmType"`
	Origin          string `json:"origin"`
	Details         string `json:"details"`
	Severity        string `json:severity`
}

type Notification struct {
	name      string
	item      InitialMonitoredServiceItem
	status    string
	timestamp int64
}
