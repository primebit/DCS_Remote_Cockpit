package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"code.google.com/p/go.net/websocket"
	"log"
	"html/template"
)

var messageQueues []chan string

func Serve(messageChan chan string) {
	go func(messageQueues *[]chan string) {
		for message := range messageChan {
			go cloneMessageToQueues(message, messageQueues)
		}
	}(&messageQueues)

	r := mux.NewRouter()
	r.HandleFunc("/", IndexPage)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)
	http.Handle("/socket", websocket.Handler(DataSocket))
	http.ListenAndServe(":8100", nil)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./views/index.html")
	if err != nil {
		log.Println("Index page:", err)
	}
	t.Execute(w, "")
}

func cloneMessageToQueues(message string, messageQueues *[]chan string) {
	for _,queue := range *messageQueues {
		if queue != nil {
			queue <- message
		}
	}
}

func DataSocket(ws *websocket.Conn) {
	queue := make(chan string)

	messageQueues = append(messageQueues, queue)
	queueId := len(messageQueues)-1

	log.Println("New websocket connection #", queueId)

	doneChan := make(chan bool, 1)
	errChan:= make(chan bool, 1)
	for message := range queue {
		if(len(doneChan) > 0) {
			continue
		}

		doneChan <- true
		go func(message string, ws *websocket.Conn, doneChan chan bool, errChan chan bool) {
			_, err := ws.Write([]byte("{\"type\":\"update\", \"data\":" + message + "}"))
			if err != nil {
				deleteQueue(queueId)
				log.Println("Websocket transmit error, connection closed #", queueId)
				errChan <- true
				return

			} else {
				response := make([]byte, 10)
				n, err := ws.Read(response)
				if err != nil || n != 2 {
					deleteQueue(queueId)
					log.Println("Websocket confirm receive error, connection closed #", queueId)
					errChan <- true
					return
				}
			}
			<- doneChan
		}(message, ws, doneChan, errChan)

		if len(errChan) > 0 {
			break
		}
	}
	ws.Close()
}

func deleteQueue(queueId int) {
	messageQueues[queueId] = nil
}
