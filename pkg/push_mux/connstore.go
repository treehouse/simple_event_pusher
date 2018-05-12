package push_mux

import (
	push "github.com/nicolasjhampton/simple_event_pusher/internal/connection"
	event "github.com/nicolasjhampton/simple_event_pusher/pkg/event"
)

type AddStruct struct {
	Key string
	Conn *push.Connection
}

type SendStruct struct {
	Key string 
	Msg *event.Message
}

type ConnStore struct {
	connList *ConnList
	AddConn chan AddStruct
	DeleteConn chan string
	SendToPush chan SendStruct
}

func New() *ConnStore {
	return &ConnStore{
		connList: newConnList(), //ConnList{},
		AddConn: make(chan AddStruct, 1),
		DeleteConn: make(chan string),
		SendToPush: make(chan SendStruct, 1),
	}
}


func (cs ConnStore) Add(channel string, c *push.Connection) {
	cs.AddConn <- AddStruct{ Key: channel, Conn: c }
}

func (cs ConnStore) Delete(channel string) {
	cs.DeleteConn <- channel
}


func (cs ConnStore) Run() {
	for {
		select {
		case incomingConn := <-cs.AddConn:
			key := incomingConn.Key//.(string)
			conn := incomingConn.Conn//.(*push.Connection)
			cs.connList.Add(key, conn)
		case key := <-cs.DeleteConn:
			cs.connList.Delete(key)
		case outgoingMsg := <-cs.SendToPush:
			key := outgoingMsg.Key//.(string)
			msg := outgoingMsg.Msg//.(*event.Message)
			cs.connList.SendToPush(key, msg)
		} 
	}
}