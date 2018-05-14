package push_mux_test

import (
	"fmt"
	mux "github.com/treehouse/simple_event_pusher/pkg/push_mux"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"reflect"
)


var pConn1 *push.Connection = push.NewConnection("channel_one")
var pConn2 *push.Connection = push.NewConnection("channel_two")
var pConn3 *push.Connection = push.NewConnection("channel_three")


func ExampleNew() {

	pushMux := mux.New()

	fmt.Println(reflect.TypeOf(pushMux))
	// Output:
	// *push_mux.ConnStore
}

func ExampleConnStore_Add() {

	pushMux := mux.New()

	go pushMux.Run()

	pushMux.Add(pConn1)
	pushMux.Add(pConn2)
	pushMux.Add(pConn3)

	channels := pushMux.ConnList.Channels()

	// fmt.Println(channels)
	fmt.Println(reflect.DeepEqual(
		channels,
		[]string{"channel_one", "channel_two", "channel_three"},
	))
	// Output:
	// true
}

func ExampleConnStore_Remove() {

	pushMux := mux.New()

	go pushMux.Run()

	pushMux.Add(pConn1)
	pushMux.Add(pConn2)
	pushMux.Add(pConn3)

	pushMux.Remove(pConn1)

	channels := pushMux.ConnList.Channels()
	
	fmt.Println(channels)

	fmt.Println(reflect.DeepEqual(
		channels,
		[]string{"channel_two", "channel_three"},
	))
	// Output:
	// true
}

func ExampleConnStore_Send() {

	pushMux := mux.New()

	go pushMux.Run()

	pushMux.Send(&event.Message{
		EventStr: "",
		Channel:  "channel_three",
		DataStr:  "{ \"a\": \"data\", \"pay\": \"load\" }",
	})

}