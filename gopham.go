/**
 * gopham
 * go push message manager
 */

package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"golang.org/x/net/websocket"
)

// Message is JSON structure
type Message struct {
	Channel string
	TTL     int
	Data    JSON
}

// JSON is json type
type JSON map[string]interface{}

// websocket connection manager
func connectionManager() (connAdd, connDel chan *websocket.Conn, connSafe func(func([]*websocket.Conn))) {
	connections := make([]*websocket.Conn, 0, 100)
	connAdd = make(chan *websocket.Conn, 1)
	connDel = make(chan *websocket.Conn, 1)
	mutex := new(sync.Mutex)

	// safety connections getter
	connSafe = func(f func([]*websocket.Conn)) {
		defer mutex.Unlock()
		mutex.Lock()
		f(connections)
	}

	// watch add event
	go func() {
		for {
			func() {
				conn := <-connAdd
				log.Println("server: new connection")
				mutex.Lock()
				defer mutex.Unlock()
				connections = append(connections, conn)
				log.Println("connections:", len(connections))
			}()
		}
	}()

	// watch delete event
	go func() {
		for {
			func() {
				conn := <-connDel
				log.Println("server: connection closed")
				mutex.Lock()
				defer mutex.Unlock()
				for i, ws := range connections {
					if ws == conn {
						connections = append(connections[:i], connections[i+1:]...)
						log.Println("connections:", len(connections))
						break
					}
				}
			}()
		}
	}()

	return
}

// websocket end point
func ws() (connSafe func(func([]*websocket.Conn))) {
	connAdd, connDel, connSafe := connectionManager()

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(
			func(ws *websocket.Conn) {
				// add connection
				connAdd <- ws

				defer func() {
					// delete connection
					connDel <- ws
				}()

				for {
					// receive message
					message := new(Message)
					websocket.JSON.Receive(ws, message)
					if message.Channel == "" && message.Data == nil {
						return
					}
					log.Printf("client: %#v\n", message)
				}
			}),
		}
		s.ServeHTTP(w, r)
	})

	return
}

func main() {
	// websocket route
	connSafe := ws()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "%s\n", "gopham works")
	})

	// post message
	engine.POST("/", func(ctx *gin.Context) {
		message := new(Message)
		err := ctx.BindJSON(message)
		if err != nil {
			ctx.JSON(400, gin.H{
				"status": "ng",
				"error":  err.Error(),
			})
			return
		}

		if message.Channel == "" || message.Data == nil {
			ctx.JSON(400, gin.H{
				"status": "ng",
				"error":  "`channel` and `data` is required.",
			})
			return
		}

		log.Printf("server: %#v\n", message)
		data := JSON{
			"channel": message.Channel,
			"ttl":     message.TTL,
			"data":    message.Data,
		}

		connectionLen := 0
		// broadcast message
		connSafe(func(connections []*websocket.Conn) {
			for _, ws := range connections {
				websocket.JSON.Send(ws, data)
			}
			connectionLen = len(connections)
		})

		ctx.JSON(200, gin.H{
			"status":      "ok",
			"connections": connectionLen,
			"message":     data,
		})
	})

	// static & middleware route
	engine.Static("/static", "static")
	engine.Use(gin.WrapH(http.DefaultServeMux))

	// listen
	log.Println("gopham server started.")
	err := engine.Run(":3000")
	if err != nil {
		panic(err)
	}
}
