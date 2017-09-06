package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"strings"
)

func connectToWebSocket(ws_channel, quite_channel chan string) {
	var origin = "http://localhost:10000"
	var url = "ws://localhost:10001/"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		quite_channel <- "Unable to dial for the websocket " + err.Error()
		return
	}

	//receive frames from websocket
	for {
		var message string
		websocket.Message.Receive(ws, &message)
		if strings.Compare(message, "") == 0 {
			//empty string received means connection closed
			quite_channel <- "Empty string received, Seems to be connection closed!"
			return
		} else {
			ws_channel <- message
		}
	}
}

func main() {
	ws_channel := make(chan string)
	quite_channel := make(chan string)
	go connectToWebSocket(ws_channel, quite_channel)

	for {
		select {
		case frame := <-ws_channel:
			fmt.Printf("Frame receieved %s \n", frame)
		case sig_quite := <-quite_channel:
			fmt.Printf("Exit signal received! %s\n", sig_quite)
			log.Fatal("Exiting since error occured")
		}
	}
}
