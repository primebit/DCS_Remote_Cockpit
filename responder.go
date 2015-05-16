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

//	origin := "http://localhost:9000/"
//	url := "ws://localhost:9000/radar/responder"
//	ws, err := websocket.Dial(url, "", origin)
//	if err != nil {
//		log.Fatal(err)
//	} else {
//		log.Println("Websocket connected")
//	}
//	if some, err := ws.Write([]byte("hello, world!\n")); err != nil {
//		log.Fatal(err)
//	} else {
//		fmt.Println("Sended: ", some)
//	}
//	log.Println("Sended")

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
