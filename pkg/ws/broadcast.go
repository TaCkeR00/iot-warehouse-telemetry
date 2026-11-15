package ws

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Clients = make(map[*Client]bool)
	Mu      sync.Mutex
)

func Broadcast(msg []byte) error {
	Mu.Lock()
	defer Mu.Unlock()
	var errs error
	for c := range Clients {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			c.Conn.Close()
			delete(Clients, c)
			errs = errors.Join(errs, err)
		}
	}
	return errs
}
