package pham

import (
	"log"
	"sync"
)

var (
	connAdd  = make(chan Connection, 1)
	connDel  = make(chan Connection, 1)
	connSafe func(func([]Connection))
)

// ConnectionManager is websocket connection manager
func ConnectionManager() (cAdd, cDel chan<- Connection, connSafe func(func([]Connection))) {
	connections := make([]Connection, 0, 100)
	cAdd = connAdd
	cDel = connDel
	connMutex := new(sync.Mutex)
	// safety connections getter
	connSafe = func(f func([]Connection)) {
		defer connMutex.Unlock()
		connMutex.Lock()
		f(connections)
	}

	// watch add & delete event
	go func() {
		for {
			func() {
				select {
				case conn := <-connAdd:
					log.Println("server: new connection")
					connMutex.Lock()
					defer connMutex.Unlock()

					connections = append(connections, conn)
					log.Println("connections:", len(connections))

				case conn := <-connDel:
					log.Println("server: connection closed")
					connMutex.Lock()
					defer connMutex.Unlock()

					for i, ws := range connections {
						if ws == conn {
							connections = append(connections[:i], connections[i+1:]...)
							log.Println("connections:", len(connections))
							break
						}
					}
				}
			}()
		}
	}()

	return
}
