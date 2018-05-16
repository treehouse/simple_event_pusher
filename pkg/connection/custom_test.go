package push_test

import (
	"fmt"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"reflect"
)

func ExampleNewCustomConn() {

	conn = push.NewCustomConn("myChannel", "*")

	fmt.Println("Returns:", reflect.TypeOf(conn))
	fmt.Println("Public properties:")
	fmt.Println("Channel:", conn.Channel())
	// Output:
	// Returns: *push.CustomConnection
	// Public properties:
	// Channel: myChannel
}

func ExampleCustomConnection_Send() {

	evt := event.Event{
		CHANNEL: "myChannel",
		EVENT:   "event",
		DATA:    "{ \"a\": \"data\", \"pay\": \"load\" }",
		ID:      "",
	}

	go conn.Send(evt)

}
