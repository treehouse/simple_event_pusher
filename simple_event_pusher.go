package main

import (
	"fmt"
	"github.com/donovanhide/eventsource"
	// "github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

var channel = "hello_world";

// Boilerplate for eventsource package
type eventPusherMessage struct {
	EventStr string `json:"event"`
	Channel  string `json:"channel"`
	DataStr  string `json:"data"`
}

func (e *eventPusherMessage) Id() string    { return "" }
func (e *eventPusherMessage) Event() string { return e.EventStr }
func (e *eventPusherMessage) Data() string  { return e.DataStr }

func main() {
	fmt.Println("Hello World")

	pusher := eventsource.NewServer()
	defer pusher.Close()

	l, err := net.Listen("tcp", ":8080")
	
	if err != nil {
			return
	}
	defer l.Close()


	http.HandleFunc("/channel", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		bcastChan := pusher.Handler(channel)
		bcastChan.ServeHTTP(w, r)
	})

	// server := &http.Server{
	// 	Addr:         ":3001",
	// 	Handler:      muxApp, // all this needs is something that fulfills the Handler interface
	// 	WriteTimeout: 15 * time.Second,
	// 	ReadTimeout:  15 * time.Second,
	// }

	msgChan := make(chan eventPusherMessage, 1)

	go http.Serve(l, nil)

	// emit a message every two seconds
	go func() {
		for {
			time.Sleep(2 * time.Second)
			msgChan <- eventPusherMessage{
				EventStr: "message",
				Channel:  channel,
				DataStr:  "{\"test\":\"message\"}",
			}
		}
	}()

	for {
		nextMsg := <-msgChan
		fmt.Println("go channel works")
		pusher.Publish([]string{nextMsg.Channel}, &nextMsg)
	}

}
