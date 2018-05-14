package client

import (
	"encoding/json"
	"fmt"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
)

type MessageInterface interface {
	GetChannel() string
	GetPayload() string
}

type Incoming interface {
	Close() error
	ReceiveMessage() (MessageInterface, error)
}

// PUBLISH event_pusher '{ "event": "", "channel": "hello-multiplex", "data": "this is data"}'
func ListenForMsgs(cl *push.ConnList, dataSource Incoming) {

	defer dataSource.Close()

	for {
		msg, err := dataSource.ReceiveMessage()

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("got payload from redis pubsub %+v\n", msg.GetPayload())

		epmsg := event.Message{}
		if err := json.Unmarshal([]byte(msg.GetPayload()), &epmsg); err != nil {
			fmt.Println(err)
			continue
		}

		go cl.Send(&epmsg)
	}
}
