package handler_test

import (
	"fmt"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	handler "github.com/treehouse/simple_event_pusher/pkg/handler"
	"reflect"
)

var connList *push.ConnList = push.NewConnList()

func ExampleServeSession() {

	pushHandler := handler.ServeSession(connList, "http://localhost:3001")

	fmt.Println(reflect.TypeOf(pushHandler))
	// Output:
	// func(http.ResponseWriter, *http.Request)
}
