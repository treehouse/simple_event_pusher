package push_test

import (
	"fmt"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"reflect"
)

var connList *push.ConnList
var pConn1 *push.Connection = push.NewConnection("channel_one")
var pConn2 *push.Connection = push.NewConnection("channel_two")
var pConn3 *push.Connection = push.NewConnection("channel_three")

func ExampleNewConnList() {

	connList = push.NewConnList()

	fmt.Println(reflect.TypeOf(connList))
	// Output:
	// *push.ConnList
}

func ExampleConnList_Add() {

	connList.Add(pConn1)
	connList.Add(pConn2)
	connList.Add(pConn3)

	channels := connList.Channels()

	fmt.Println(reflect.DeepEqual(
		channels,
		[]string{"channel_one", "channel_two", "channel_three"},
	))
	// Output:
	// true
}

func ExampleConnList_Remove() {

	connList.Remove(pConn2)

	channels := connList.Channels()

	fmt.Println(reflect.DeepEqual(
		channels,
		[]string{"channel_one", "channel_three"},
	))
	// Output:
	// true
}

func ExampleConnList_Send() {

	connList.Send(&event.Message{
		EventStr: "",
		Channel:  "channel_three",
		DataStr:  "{ \"a\": \"data\", \"pay\": \"load\" }",
	})

}
