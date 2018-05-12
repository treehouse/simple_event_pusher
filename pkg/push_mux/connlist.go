package push_mux

import (
	"sync"
	push "github.com/nicolasjhampton/simple_event_pusher/internal/connection"
	event "github.com/nicolasjhampton/simple_event_pusher/pkg/event"
)

type ConnMap map[string]*push.Connection

type ConnList struct {
	list ConnMap
	mu sync.RWMutex
}

func newConnList() *ConnList {
	return &ConnList{
		list: ConnMap{},
		mu: sync.RWMutex{},
	}
}

func (cl *ConnList) Add(sessionName string, conn *push.Connection) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.list[sessionName] = conn
}

func (cl *ConnList) Delete(sessionName string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	delete(cl.list, sessionName)
}

func (cl *ConnList) SendToPush(sessionName string, msg *event.Message) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	if conn, ok := cl.list[sessionName]; ok {
    conn.SendToPush(msg)
	}
}