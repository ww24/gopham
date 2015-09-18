package pham

import "github.com/gin-gonic/gin"

// ServerSentEventsConnection structure
type ServerSentEventsConnection struct {
	ctx *gin.Context
}

// Send implemented Connection interface
func (sse *ServerSentEventsConnection) Send(data JSON) (err error) {
	sse.ctx.SSEvent("message", data)
	sse.ctx.Writer.Flush()
	return
}

// SSEHandler for gin framework route handler
func SSEHandler(ctx *gin.Context) {
	connection := &ServerSentEventsConnection{ctx: ctx}
	// add connection
	connAdd <- connection

	defer func() {
		// delete connection
		connDel <- connection
	}()

	// set sse header
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Writer.Flush()

	// watch close connection event
	closer := ctx.Writer.CloseNotify()
	for {
		select {
		case <-closer:
			return
		}
	}
}
