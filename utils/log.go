package utils

import (
	"log"
	"os"
)

var (
	Log *log.Logger
)

func NewLog(logpath string) {
	println("LogFile: " + logpath)
	file, err := os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}
