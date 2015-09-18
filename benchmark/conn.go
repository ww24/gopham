/**
 * go run conn.go > conn.log 2>&1
 * sudo sysctl -w kern.maxfiles=65536
 * ulimit -n 65536
 */

package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	log.Println("WebSocket Bench")
	wrapper(func() (ch chan bool) {
		ch = wsCli("ws://127.0.0.1:3000/subscribe")
		return
	})

	log.Println("Server-Sent Events Bench")
	wrapper(func() (ch chan bool) {
		ch = sseCli("http://127.0.0.1:3000/sse")
		return
	})
}

func wrapper(client func() chan bool) {
	chs := make([]chan bool, 10000)

	for i := range chs {
		log.Println(i)
		chs[i] = client()
	}

	jsonBuffer := bytes.NewBuffer([]byte(`{
		"channel": "test",
		"ttl": 0,
		"data": {"message": "json"}
	}`))
	res, err := http.Post("http://127.0.0.1:3000/", "application/json", jsonBuffer)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))

	for _, ch := range chs {
		<-ch
	}
}

func sseCli(url string) (ch chan bool) {
	defer catchError()

	ch = make(chan bool, 1)

	res, err := http.Get(url)
	if err != nil {
		go func() {
			ch <- true
		}()
		log.Println("sse: connection failure")
		panic(err)
	}

	sc := bufio.NewScanner(res.Body)

	go func() {
		defer res.Body.Close()
		defer func() {
			ch <- true
		}()

		for {
			switch {
			case res.Close:
				log.Println("closed")
				return
			case sc.Scan():
				line := sc.Text()
				log.Println(line)

				// auto close
				if line == "\n" {
					return
				}
			}
		}
	}()

	return
}

func wsCli(url string) (ch chan bool) {
	defer catchError()

	ch = make(chan bool, 1)

	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		go func() {
			ch <- true
		}()
		log.Println("ws: connection failure")
		panic(err)
	}

	go func() {
		defer ws.Close()
		defer func() {
			ch <- true
		}()

		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				if err == io.EOF {
					log.Println("close event")
					return
				}
				panic(err)
			}
			log.Println(msg)

			// auto close
			return
		}
	}()

	return
}

func catchError() {
	cause := recover()
	if cause == nil {
		return
	}

	err, check := cause.(error)
	if check == false {
		panic(cause)
	}

	log.Println("error:", err.Error())
}
