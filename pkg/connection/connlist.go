package push

import (
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"sync"
)

// Base data structure storing connections. Wrapped with ConnList
// to ensure threadsafe reads/writes.
type ConnMap map[string]*Connection

// Essencially a light wrapper around a map of connections with
// methods and a Read/Write Mutex to ensure threadsafe reads and
// writes. Read operations are allowed to occur concurrently while
// locking out writes. Write operations lock for exclusive access.
type ConnList struct {
	list ConnMap
	mu   sync.Mutex
}

// Creates and initializes a new ConnList
func NewConnList() *ConnList {
	return &ConnList{
		list: ConnMap{},
		mu:   sync.Mutex{},
	}
}

// Adds a new connection to the ConnList, locking out all other
// threads while write occurs.
func (cl *ConnList) Add(conn *Connection) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.list[conn.Channel] = conn
}

// Removes a connection to the ConnList, locking out all other
// threads while write occurs.
func (cl *ConnList) Remove(conn *Connection) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	delete(cl.list, conn.Channel)
}

// Pushes a message to connected browser. Locks write operations
// out while push occurs, but allows for concurrent pushes by
// multiple threads. Reads channel from message to determine
// correct browser to send to.
func (cl *ConnList) SendToPush(msg *event.Message) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if conn, ok := cl.list[msg.Channel]; ok {
		conn.SendToPush(msg)
	}
}

// Helper method primarily for checking private state of channel
// list during testing.
func (cl *ConnList) Channels() []string {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	channels := []string{}
	for channel, _ := range cl.list {
		channels = append(channels, channel)
	}
	return channels
}
