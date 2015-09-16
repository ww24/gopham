package pham

import (
	"log"
	"sync"
)

var (
	connAdd = make(chan Connection, 1)
	connDel = make(chan Connection, 1)
)

// Connection interface
type Connection interface {
	Send(data JSON) (err error)
}

// ConnectionManager is websocket connection manager
func ConnectionManager() (cAdd, cDel chan Connection, connSafe func(func([]Connection))) {
	connections := make([]Connection, 0, 100)
	cAdd = connAdd
	cDel = connDel
	mutex := new(sync.Mutex)

	// safety connections getter
	connSafe = func(f func([]Connection)) {
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
