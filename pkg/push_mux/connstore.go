package push_mux

import (
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
)

// A wrapper around the threadsafe connList.
type ConnStore struct {
	ConnList   *push.ConnList
	AddConn    chan *push.Connection
	RemoveConn chan *push.Connection
	SendToPush chan *event.Message
	requestChannels chan int
	channels chan []string
}

// Creates a new ConnStore. 
func New() *ConnStore {
	return &ConnStore{
		ConnList:   push.NewConnList(),
		AddConn:    make(chan *push.Connection),
		RemoveConn: make(chan *push.Connection),
		SendToPush: make(chan *event.Message),
		requestChannels: make(chan int),
		channels: make(chan []string),
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

func (cs ConnStore) Channels() []string {
	cs.requestChannels <- 1
	channels := <- cs.channels
	return channels
}

// Run executes read/write requests to the underlying, threadsafe 
// connList. Meant to be ran as a separate, continuously running 
// goroutine. I think I need to reduce this down to one channel 
// order to get the right timing. Thought I had it.
func (cs ConnStore) Run() {
	for {
		select {
		case newConn := <-cs.AddConn:
			cs.ConnList.Add(newConn)
		case conn := <-cs.RemoveConn:
			cs.ConnList.Remove(conn)
		case msg := <-cs.SendToPush:
			cs.ConnList.SendToPush(msg)
		case <-cs.requestChannels:
			channels := cs.ConnList.Channels()
			cs.channels <- channels
		}
	}
}
