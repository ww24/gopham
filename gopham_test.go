package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ww24/gopham/pham"
	"github.com/ww24/gopham/pham/client"

	"golang.org/x/net/websocket"
)

func post(url string, data pham.JSON) (str string, err error) {
	jsonb, err := json.Marshal(data)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(jsonb)
	res, err := http.Post(url, "application/json", buf)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	str = string(body)

	return
}

func TestWebSocket(t *testing.T) {
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	ws, err := websocket.Dial("ws://"+ts.Listener.Addr().String()+"/subscribe", "", "http://localhost/")
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	ch := make(chan bool, 1)
	go func() {
		msg := make([]byte, 512)
		size, err := ws.Read(msg)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("TestWebSocket:received:", size, string(msg[:size]))

		ch <- true
	}()

	// send realtime message
	res, err := post(ts.URL, pham.JSON{
		"channel": "test",
		"data": pham.JSON{
			"type": "ping",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println("TestWebSocket:sent:", res)

	<-ch
}

func TestServerSentEvents(t *testing.T) {
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	cli, err := client.NewSSEClient(ts.URL + "/sse")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ch := make(chan bool, 1)
	go func() {
		data := <-cli.Listener
		log.Println("TestServerSentEvents:received:", data)

		ch <- true
	}()

	// send realtime message
	res, err := post(ts.URL, pham.JSON{
		"channel": `test`,
		"data": pham.JSON{
			"type": "ping",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println("TestServerSentEvents:sent:", res)

	<-ch
}

func BenchmarkWebSocket(b *testing.B) {
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	wg := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		ws, err := websocket.Dial("ws://"+ts.Listener.Addr().String()+"/subscribe", "", "http://localhost/")
		if err != nil {
			panic(err)
		}
		defer ws.Close()

		wg.Add(1)

		go func() {
			msg := make([]byte, 512)
			size, err := ws.Read(msg)
			if err != nil {
				log.Println("err:", err)
			}
			log.Println("TestWebSocket:received:", size, string(msg[:size]))

			wg.Done()
		}()
	}

	// send realtime message
	res, err := post(ts.URL, pham.JSON{
		"channel": "test",
		"data": pham.JSON{
			"type": "ping",
		},
	})
	if err != nil {
		panic(err)
	}
	log.Println("TestWebSocket:sent:", res)

	wg.Wait()
}

func BenchmarkServerSentEvents(b *testing.B) {
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	wg := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		cli, err := client.NewSSEClient(ts.URL + "/sse")
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		wg.Add(1)

		go func() {
			data := <-cli.Listener
			log.Println("TestServerSentEvents:received:", data)

			wg.Done()
		}()
	}

	// send realtime message
	res, err := post(ts.URL, pham.JSON{
		"channel": `test`,
		"data": pham.JSON{
			"type": "ping",
		},
	})
	if err != nil {
		panic(err)
	}
	log.Println("TestServerSentEvents:sent:", res)

	wg.Wait()
}
