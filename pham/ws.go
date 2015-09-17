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
		w.WriteHeader(101)
		s.ServeHTTP(w, r)
	})
}
