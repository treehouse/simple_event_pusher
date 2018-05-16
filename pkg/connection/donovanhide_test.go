package push_test

import (
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"fmt"
	"reflect"
)

func ExampleNewDonovanConn() {

	conn = push.NewDonovanConn("myChannel")

	fmt.Println("Returns:", reflect.TypeOf(conn))
	fmt.Println("Public properties:")
	fmt.Println("Channel:", conn.Channel())
	// Output:
	// Returns: *push.DonovanhideConnection
	// Public properties:
	// Channel: myChannel
}

func ExampleDonovanhideConnection_Send() {

	evt := event.Event{
		CHANNEL: "myChannel",
		EVENT:   "event",
		DATA:    "{ \"a\": \"data\", \"pay\": \"load\" }",
		ID:      "",
	}

	go conn.Send(evt)

}

func ExampleDonovanhideConnection_Msgs() {

	// closes push connection once outside function finishes.
	// defer conn.Close()
	go conn.Msgs()

	// ServePUSH acts as a blocking line here. While blocking, Msgs
	// continues to run on it's own thread, listening for msgs to push.
	// When the browser disconnects, ServePUSH stops running, and the
	// outer function reaches it's end, triggering the deferred
	// connection close and ending the Msgs goroutine.

	// conn.ServePUSH(w, r)
}
