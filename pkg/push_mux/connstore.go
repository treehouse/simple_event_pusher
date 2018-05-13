package push_mux

import (
	push "github.com/treehouse/simple_event_pusher/internal/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
)

type AddStruct struct {
	Key  string
	Conn *push.Connection
}

// type SendStruct struct {
// 	Key string
// 	Msg *event.Message
// }

type ConnStore struct {
	connList   *push.ConnList
	AddConn    chan AddStruct
	DeleteConn chan string
	SendToPush chan *event.Message
}

func New() *ConnStore {
	return &ConnStore{
		connList:   push.NewConnList(),
		AddConn:    make(chan AddStruct, 1),
		DeleteConn: make(chan string),
		SendToPush: make(chan *event.Message, 1),
	}
}

func (cs ConnStore) Add(channel string, c *push.Connection) {
	cs.AddConn <- AddStruct{Key: channel, Conn: c}
}

func (cs ConnStore) Delete(channel string) {
	cs.DeleteConn <- channel
}

func (cs ConnStore) Send(msg *event.Message) {
	cs.SendToPush <- msg
}

func (cs ConnStore) Run() {
	for {
		select {
		case incomingConn := <-cs.AddConn:
			key := incomingConn.Key
			conn := incomingConn.Conn
			cs.connList.Add(key, conn)
		case key := <-cs.DeleteConn:
			cs.connList.Delete(key)
		case outgoingMsg := <-cs.SendToPush:
			// key := outgoingMsg.Channel
			msg := outgoingMsg
			cs.connList.SendToPush(msg)
		}
	}
}
