package push_mux

import (
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
)

// A wrapper around the threadsafe connList that, in theory, can be 
// adjusted to create a rather large, fast buffer of msgs and 
// connections in each channel, and many spawned concurrent 
// goroutines. The channels buffer size can be increased, and 
// goroutines can continue to spawn if connList is blocking. Given
// that the majority of actions used are "SendToPush", potentially
// access could be blocked for several thousand "SendToPush" goroutines,
// but all of them can access the connList concurrently as soon as 
// the lock is cleared.
type ConnStore struct {
	connList   *push.ConnList
	AddConn    chan *push.Connection
	RemoveConn chan *push.Connection
	SendToPush chan *event.Message
}

// Creates a new ConnStore. Channels on the ConnStore can be buffered. 
// Provide the buffer length as the second argument to make to 
// initialize a buffered channel:
//
//		&ConnStore{
//				connList:   push.NewConnList(),
//				AddConn:    make(chan *push.Connection, 100),
//				RemoveConn: make(chan *push.Connection, 100),
//				SendToPush: make(chan *event.Message, 1000),
//		}
//
// Sends to a buffered channel block only when the buffer is full.
// Receives block when the buffer is empty.
// TODO: Modify this code with the appropreate buffer for our needs.
func New() *ConnStore {
	return &ConnStore{
		connList:   push.NewConnList(),
		AddConn:    make(chan *push.Connection, 1),
		RemoveConn: make(chan *push.Connection, 1),
		SendToPush: make(chan *event.Message, 1),
	}
}

func (cs ConnStore) Add(c *push.Connection) {
	cs.AddConn <- c
}

func (cs ConnStore) Remove(c *push.Connection) {
	cs.RemoveConn <- c
}

func (cs ConnStore) Send(msg *event.Message) {
	cs.SendToPush <- msg
}

// Run builds up and executes read/write requests to the underlying,
// threadsafe connList. Meant to be ran as a separate, continuously 
// running goroutine.
//
// Thinking about running each of these connList functions
// as goroutines to allow faster, safe access to connList
// in situations where the channel buffers are backed up.
// https://stackoverflow.com/questions/8509152/max-number-of-goroutines?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa
func (cs ConnStore) Run() {
	for {
		select {
		case newConn := <-cs.AddConn:
			cs.connList.Add(newConn)
		case conn := <-cs.RemoveConn:
			cs.connList.Remove(conn)
		case outgoingMsg := <-cs.SendToPush:
			msg := outgoingMsg
			cs.connList.SendToPush(msg)
		}
	}
}
