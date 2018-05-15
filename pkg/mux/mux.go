package mux

import (
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"sync"
	"fmt"
)

// Base data structure storing connections. Wrapped with ConnList
// to ensure threadsafe reads/writes.
type ConnMap map[string]push.Connection

// Essencially a light wrapper around a map of connections with
// methods and a Read/Write Mutex to ensure threadsafe reads and
// writes. May be able to improve performance with a RWMutex
// that allows concurrent reads while locking writes.
type ConnStore struct {
	list ConnMap
	mu   sync.Mutex
}

// Creates and initializes a new ConnList
func New() *ConnStore {
	return &ConnStore{
		list: ConnMap{},
		mu:   sync.Mutex{},
	}
}

// Adds a new connection to the ConnList, locking out all other
// threads while write occurs.
func (cs *ConnStore) Add(conn push.Connection) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.list[conn.Channel()] = conn
}

// Removes a connection to the ConnList, locking out all other
// threads while write occurs.
func (cs *ConnStore) Remove(conn push.Connection) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.list, conn.Channel())
}

// Pushes a message to connected browser, locking out all other
// threads while write occurs. Reads channel from message to determine
// correct browser to send to.
func (cs *ConnStore) Send(msg event.Message) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	fmt.Println(msg)
	if conn, ok := cs.list[msg.GetChannel()]; ok {
		conn.Send(msg)
	}
}

// Helper method primarily for checking private state of channel
// list during testing.
func (cs *ConnStore) Channels() []string {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	channels := []string{}
	for channel, _ := range cs.list {
		channels = append(channels, channel)
	}
	return channels
}
