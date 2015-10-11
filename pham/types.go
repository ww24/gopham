package pham

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
