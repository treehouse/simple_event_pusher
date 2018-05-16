package client

import (
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	mux "github.com/treehouse/simple_event_pusher/pkg/mux"
	"fmt"
)

type Incoming interface {
	Close() error
	ReceiveMessage() (event.Message, error)
}

func ListenForMsgs(mux *mux.ConnStore, dataSource Incoming) {

	defer dataSource.Close()

	for {
		msg, err := dataSource.ReceiveMessage()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go mux.Send(msg)
	}
}
