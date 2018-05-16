package handler_test

import (
	handler "github.com/treehouse/simple_event_pusher/pkg/handler"
	mux "github.com/treehouse/simple_event_pusher/pkg/mux"
	"fmt"
	"reflect"
)

var connList *mux.ConnStore = mux.New()

func ExampleServeSession() {

	pushHandler := handler.ServeSession(connList, "http://localhost:3001")

	fmt.Println(reflect.TypeOf(pushHandler))
	// Output:
	// func(http.ResponseWriter, *http.Request)
}
