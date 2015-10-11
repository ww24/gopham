package pham

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// WebSocketConnection structure
type WebSocketConnection struct {
	ws *websocket.Conn
}

// Send implemented Connection interface
func (wc *WebSocketConnection) Send(data []byte) (err error) {
	_, err = wc.ws.Write(data)
	return
}

// WSHandler is websocket end point
func WSHandler(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{Handler: websocket.Handler(
		func(ws *websocket.Conn) {
			connection := &WebSocketConnection{ws: ws}

			// add connection
			connAdd <- connection

			defer func() {
				// delete connection
				connDel <- connection
			}()

			for {
				// receive message
				message := new(Message)
				err := websocket.JSON.Receive(ws, message)
				if err != nil {
					// close event
					if err == io.EOF {
						return
					}

					log.Println(err)
				}
				log.Printf("client: %#v\n", message)
			}
		}),
	}
	s.ServeHTTP(w, r)
}
