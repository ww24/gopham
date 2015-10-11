package pham

import "sync"

// Message is JSON structure
type Message struct {
	Channel string
	TTL     int
	Data    JSON
}

// JSON is json type
type JSON map[string]interface{}

// Connection interface
type Connection interface {
	Send(data []byte) (err error)
}

// ConnectionManager structure
type ConnectionManager struct {
	connections []Connection
	connMutex   *sync.Mutex
	connAdd     chan Connection
	connDel     chan Connection
}
