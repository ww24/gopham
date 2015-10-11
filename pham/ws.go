package pham

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// webSocketConnection structure
type webSocketConnection struct {
	ws *websocket.Conn
}

// Send implemented Connection interface
func (wc *webSocketConnection) Send(data []byte) (err error) {
	_, err = wc.ws.Write(data)
	return
}

// WSHandler is websocket end point
func WSHandler(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{Handler: websocket.Handler(
		func(ws *websocket.Conn) {
			connection := &webSocketConnection{ws: ws}

			// add connection
			manager.AddConnection(connection)
			defer func() {
				// delete connection
				manager.DelConnection(connection)
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
