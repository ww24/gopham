/**
 * gopham
 * go push message manager
 */

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ww24/gopham/pham"
)

var (
	manager = pham.NewConnectionManager()
)

func main() {
	port := flag.Int("port", 3000, "Set port number.")
	modes := []string{gin.DebugMode, gin.ReleaseMode, gin.TestMode}
	mode := flag.String("mode", gin.ReleaseMode, "Set Gin Web Framework mode. ["+strings.Join(modes, " or ")+"]")
	flag.Parse()

	gin.SetMode(*mode)
	listener, errch := Serve(":"+strconv.Itoa(*port), NewHandler())

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT)
		for {
			select {
			case s := <-sig:
				log.Println(s)
				listener.Close()
			}
		}
	}()

	log.Println("gopham server started at", listener.Addr())
	log.Println("error:", <-errch)
}

// Serve is server bootstrap
func Serve(addr string, router http.Handler) (listener net.Listener, errch <-chan error) {
	ch := make(chan error)
	errch = ch

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		ch <- err
	}

	go func() {
		ch <- http.Serve(listener, router)
	}()

	return
}

// NewHandler is return http.Handler
func NewHandler() http.Handler {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "%s\n", "gopham works")
	})

	// Server-Sent Events
	router.GET("/sse", gin.WrapF(pham.SSEHandler))
	// WebSocket
	router.GET("/subscribe", gin.WrapF(pham.WSHandler))

	// post message
	router.POST("/", func(ctx *gin.Context) {
		defer func() {
			cause := recover()
			if cause != nil {
				ctx.JSON(400, gin.H{
					"status": "ng",
					"error":  cause.(error).Error(),
				})
			}
		}()

		message := new(pham.Message)
		err := ctx.BindJSON(message)
		if err != nil {
			panic(err)
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

		// encode json
		connectionLen, err := manager.Broadcast(data)
		if err != nil {
			panic(err)
		}

		ctx.JSON(200, gin.H{
			"status":      "ok",
			"connections": connectionLen,
			"message":     data,
		})
	})

	// static & middleware route
	router.Static("/static", "static")

	return router
}
