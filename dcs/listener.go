package dcs

import (
	"log"
	"net"
	"bufio"
	"sync"
)

type Message struct {
	time float64
	user string
	altBar float64
	altRad float64
	trueSpeed float64
	verticalSpeed float64
	heading	float64
}

type MassageHandler func(string, *sync.WaitGroup)

func ListenPort(port int, handleFunc MassageHandler) {
	//portStr := strings.Join([]string{":", string(port)}, "")
	portStr := ":9514"
	ln, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Opened port 9514")
	}

	// run loop forever (or until ctrl-c)
	for {
		// accept connection on port
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Accepted incoming connection")
		}

		go handleRequest(conn, handleFunc)
	}
}

func handleRequest(conn net.Conn, handleFunc MassageHandler) {
	var wg sync.WaitGroup
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		// output message received
		if err == nil && message != "" {
			wg.Add(1)
			go handleFunc(message, &wg)
		}
	}
	wg.Wait()
}
