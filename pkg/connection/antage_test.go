package push_test

import (
	"fmt"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	// cli "github.com/treehouse/simple_event_pusher/pkg/client"
	"reflect"
)

var conn push.Connection

func ExampleNewAntageConn() {

	conn = push.NewAntageConn("myChannel")

	fmt.Println("Returns:", reflect.TypeOf(conn))
	fmt.Println("Public properties:")
	fmt.Println("Channel:", conn.Channel())
	// Output:
	// Returns: *push.AntageConnection
	// Public properties:
	// Channel: myChannel
}

func ExampleConnection_Send() {

	// msg := cli.NewRedisMsg(
	// 	"one_channel_per_broswer_connection",
	// 	"{ \"a\": \"data\", \"pay\": \"load\" }",
	// )

	evt := event.Event{
		CHANNEL: "myChannel",
		EVENT:   "event",
		DATA:    "{ \"a\": \"data\", \"pay\": \"load\" }",
		ID:      "",
	}

	go conn.Send(evt)

}

func ExampleConnection_Msgs() {

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
