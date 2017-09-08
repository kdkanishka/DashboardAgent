package main

type Command interface {
	exec()
}

type PublishAlarm struct {
	notifications_map map[string]Notification
}

type ResetAlarm struct {
	notification Notification
}

type DoNothing struct {
}

func (command PublishAlarm) exec() {
}

func (command ResetAlarm) exec() {
}

func (Command DoNothing) exec() {
}
