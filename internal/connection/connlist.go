package push

import (
	event "github.com/nicolasjhampton/simple_event_pusher/pkg/event"
	"sync"
)

type ConnMap map[string]*Connection

type ConnList struct {
	list ConnMap
	mu   sync.RWMutex
}

func NewConnList() *ConnList {
	return &ConnList{
		list: ConnMap{},
		mu:   sync.RWMutex{},
	}
}

func (cl *ConnList) Add(sessionName string, conn *Connection) {
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
