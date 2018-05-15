package client

import (
	"fmt"
	// push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	mux "github.com/treehouse/simple_event_pusher/pkg/mux"
)

type Incoming interface {
	Close() error
	ReceiveMessage() (event.Message, error)
}

// PUBLISH event_pusher '{ "event": "", "channel": "hello-multiplex", "data": "this is data"}'
func ListenForMsgs(mux *mux.ConnStore, dataSource Incoming) {

	defer dataSource.Close()

	for {
		msg, err := dataSource.ReceiveMessage()

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("got payload from redis pubsub %+v\n", msg.Data())

		/// channel := msg.GetChannel() // [len(s.redisPubsubPrefix)+1:] myers has a two-part channel name here with a slice
		///channel := msg.GetChannel()[len("event_pusher") + 1:]
		// evt := event.Event{
		// 	CHANNEL: msg.Channel(),
		// 	Event: msg.GetEvent(),
		// 	Data: msg.GetData(),
		// 	Id: msg.GetId(),
		// }

		go mux.Send(msg)
	}
}
