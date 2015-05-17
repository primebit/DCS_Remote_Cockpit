package main

import (
	"log"
	"sync"
	//"golang.org/x/net/websocket"

	"dcs/responder/dcs"
	"dcs/responder/server"
)

var messageBus chan string

func main() {
	messageBus = make(chan string)
    log.Println("Launching server...")

	go dcs.ListenPort(9514, MessageHandler)
	go server.Serve(messageBus)

	select {}
}

func MessageHandler(message string, wg *sync.WaitGroup) {
	go func(message string) {
		messageBus <- message
	}(message)
	defer wg.Done()
}
