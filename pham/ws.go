package pham

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// WebSocketConnection structure
type WebSocketConnection struct {
	ws *websocket.Conn
}

// Send implemented Connection interface
func (wc *WebSocketConnection) Send(data JSON) (err error) {
	err = websocket.JSON.Send(wc.ws, data)
	return
}

// WS is websocket end point
func WS() {
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
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
					websocket.JSON.Receive(ws, message)
					if message.Channel == "" && message.TTL == 0 && message.Data == nil {
						return
					}
					log.Printf("client: %#v\n", message)
				}
			}),
		}
		s.ServeHTTP(w, r)
	})
}
