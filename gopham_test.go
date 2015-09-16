package main

import (
	"log"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestWebSocket(t *testing.T) {
	go main()
	time.Sleep(time.Second * 2)

	ws, err := websocket.Dial("ws://localhost:3000/subscribe", "", "http://localhost/")
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	cli := client()

	ch := make(chan bool)
	go func() {
		msg := make([]byte, 512)
		size, err := ws.Read(msg)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("test1:", size, string(msg[:size]))

		ch <- true
	}()

	send(ws)
	send(ws)

	<-ch
	<-cli
}

func send(ws *websocket.Conn) {
	msg := JSON{
		"channel": "test",
		"data": JSON{
			"type": "ping",
		},
	}
	err := websocket.JSON.Send(ws, msg)
	if err != nil {
		panic(err)
	}
}

func client() (ch chan bool) {
	ch = make(chan bool, 1)

	ws, err := websocket.Dial("ws://localhost:3000/subscribe", "", "http://localhost/")
	if err != nil {
		panic(err)
	}

	go func() {
		defer ws.Close()

		msg := make([]byte, 512)
		size, err := ws.Read(msg)
		if err != nil {
			panic(err)
		}
		log.Println("test2:", size, string(msg[:size]))

		ch <- true
	}()

	send(ws)
	send(ws)

	return
}
