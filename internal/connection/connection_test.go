package push_test

import (
	"fmt"
	"reflect"
	push "github.com/treehouse/simple_event_pusher/internal/connection"
	event "github.com/treehouse/simple_event_pusher/internal/event"
)

var conn *push.Connection

func ExampleNewConnection() {

	conn = push.NewConnection()

	connType := reflect.TypeOf(conn)
	fmt.Println("Returns", connType)
	// Output: 
	// Returns *push.Connection
}

func ExampleConnection_Handler() {

	handler := conn.Handler("session-channel")

	handlerType := reflect.TypeOf(handler)
	fmt.Println("Returns", handlerType)
	// Output: 
	// Returns http.HandlerFunc
}

func ExampleConnection_SendToPush() {
	
	go conn.SendToPush(&event.Message{
		EventStr: "events_can_be_assigned_to_different_handlers_on_client",
		Channel: "one_channel_per_broswer_connection",
		DataStr: "{ \"a\": \"data\", \"pay\": \"load\" }",
	})

}

func ExampleConnection_Msgs() {
	
	// closes push connection once outside function finishes.
	defer conn.Close()
	go conn.Msgs()

	// This Handler is a blocking line that stops running when browser 
	// disconnects. While blocking, Msgs continues to run, pushing msgs.
	// http.HandlerFunc.ServeHTTP(w, r)
}

