package push_mux_test

import (
	"fmt"
	mux "github.com/treehouse/simple_event_pusher/pkg/push_mux"
	// push "github.com/treehouse/simple_event_pusher/pkg/connection"
	// event "github.com/treehouse/simple_event_pusher/pkg/event"
	"reflect"
)

var pushMux *mux.ConnStore

func ExampleNew() {

	pushMux = mux.New()

	fmt.Println(reflect.TypeOf(pushMux))
	// Output:
	// *push_mux.ConnStore
}

func ExampleConnStore_Add() {
	pushMux.Add()
}

func ExampleConnStore_Remove() {
	pushMux.Remove()
}

func ExampleConnStore_SendToPush() {
	pushMux.SendToPush()
}