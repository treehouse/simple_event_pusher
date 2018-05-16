package mux_test

import (
	mux "github.com/treehouse/simple_event_pusher/pkg/mux"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	"fmt"
	"reflect"
	"sort"
)

var pmux *mux.ConnStore
var pConn1 push.Connection = push.NewAntageConn("channel_1")
var pConn2 push.Connection = push.NewAntageConn("channel_2")
var pConn3 push.Connection = push.NewAntageConn("channel_3")

func ExampleNewConnStore() {

	pmux = mux.New()

	fmt.Println(reflect.TypeOf(pmux))
	// Output:
	// *mux.ConnStore
}

func ExampleConnStore_Add() {

	pmux.Add(pConn1)
	pmux.Add(pConn2)
	pmux.Add(pConn3)

	channels := pmux.Channels()
	sort.Strings(channels)

	fmt.Println(channels)

	// Output:
	// [channel_1 channel_2 channel_3]
}

func ExampleConnStore_Remove() {

	pmux.Remove(pConn2)

	channels := pmux.Channels()
	sort.Strings(channels)

	fmt.Println(channels)

	// Output:
	// [channel_1 channel_3]
}

func ExampleConnStore_Send() {

	evt := event.Event{
		CHANNEL: "myChannel",
		EVENT:   "event",
		DATA:    "{ \"a\": \"data\", \"pay\": \"load\" }",
		ID:      "",
	}

	pmux.Send(evt)

}
