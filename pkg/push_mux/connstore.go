package push_mux

import (
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
)

// type AddStruct struct {
// 	Key  string
// 	Conn *push.Connection
// }

type ConnStore struct {
	connList   *push.ConnList
	AddConn    chan *push.Connection
	RemoveConn chan *push.Connection
	SendToPush chan *event.Message
}

// Channels can be buffered. Provide the buffer length as the second
// argument to make to initialize a buffered channel:

// ch := make(chan int, 100)
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
