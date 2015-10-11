package pham

import (
	"io"
	"net/http"
)

// ServerSentEventsConnection structure
type ServerSentEventsConnection struct {
	w *http.ResponseWriter
}

// Send implemented Connection interface
func (sse *ServerSentEventsConnection) Send(data []byte) (err error) {
	w := *sse.w

	// send data
	io.WriteString(w, `event: message
data: `+string(data)+"\n\n")

	// flush data
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	return
}

// SSEHandler for gin framework route handler
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	connection := &ServerSentEventsConnection{w: &w}

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

	// add connection
	connAdd <- connection
	defer func() {
		// delete connection
		connDel <- connection
	}()

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
