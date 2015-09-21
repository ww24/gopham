package pham

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// ServerSentEventsConnection structure
type ServerSentEventsConnection struct {
	w *http.ResponseWriter
}

// Send implemented Connection interface
func (sse *ServerSentEventsConnection) Send(data JSON) (err error) {
	w := *sse.w

	// encode json
	bytes, err := json.Marshal(data)
	if err != nil {
		io.WriteString(w, "data: {\"status\": \"ng\", \"error\":"+err.Error()+"}\n")
		return
	}

	// send data
	io.WriteString(w, `event: message
data: `+string(bytes)+"\n\n")

	// flush data
	if flusher, ok := w.(http.Flusher); ok {
		log.Println("flush")
		flusher.Flush()
	}

	return
}

// SSEHandler for gin framework route handler
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	connection := &ServerSentEventsConnection{w: &w}
	// add connection
	connAdd <- connection

	defer func() {
		// delete connection
		connDel <- connection
	}()

	// set sse header
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	header.Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)

	// flush header
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// watch close connection event
	var closer <-chan bool
	if notifier, ok := w.(http.CloseNotifier); ok {
		closer = notifier.CloseNotify()
	}

	for {
		select {
		case <-closer:
			return
		}
	}
}
