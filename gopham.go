/**
 * gopham
 * go push message manager
 */

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ww24/gopham/pham"
)

var (
	connAdd, connDel, connSafe = pham.ConnectionManager()
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "%s\n", "gopham works")
	})

	// Server-Sent Events
	engine.GET("/sse", gin.WrapF(pham.SSEHandler))
	// WebSocket
	engine.GET("/subscribe", gin.WrapF(pham.WSHandler))

	// post message
	engine.POST("/", func(ctx *gin.Context) {
		message := new(pham.Message)
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
		data := pham.JSON{
			"channel": message.Channel,
			"ttl":     message.TTL,
			"data":    message.Data,
		}

		connectionLen := 0
		// broadcast message
		connSafe(func(connections []pham.Connection) {
			for _, connection := range connections {
				connection.Send(data)
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

	// listen
	log.Println("gopham server started.")
	err := engine.Run(":3000")
	if err != nil {
		panic(err)
	}
}
