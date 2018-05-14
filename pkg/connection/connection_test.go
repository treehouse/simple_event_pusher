package push_test

import (
	"fmt"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"reflect"
)

var conn *push.Connection

func ExampleNewConnection() {

	conn = push.NewConnection("myChannel")

	fmt.Println("Returns:", reflect.TypeOf(conn))
	fmt.Println("Public properties:")
	fmt.Println("Channel:", conn.Channel)
	// Output:
	// Returns: *push.Connection
	// Public properties:
	// Channel: myChannel
}

func ExampleConnection_SendToPush() {

	go conn.SendToPush(&event.Message{
		EventStr: "events_can_be_assigned_to_different_handlers_on_client",
		Channel:  "one_channel_per_broswer_connection",
		DataStr:  "{ \"a\": \"data\", \"pay\": \"load\" }",
	})

}

func ExampleConnection_Msgs() {

	// closes push connection once outside function finishes.
	defer conn.Close()
	go conn.Msgs()

	// ServePUSH acts as a blocking line here. While blocking, Msgs
	// continues to run on it's own thread, listening for msgs to push.
	// When the browser disconnects, ServeHTTP stops running, and the
	// outer function reaches it's end, triggering the deferred
	// connection close and ending the Msgs goroutine.

	// conn.ServePUSH(w, r)
}
