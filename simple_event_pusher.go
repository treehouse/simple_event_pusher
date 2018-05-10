package main

import (
	// "bufio"
	"fmt"
	"github.com/donovanhide/eventsource"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

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
	// muxApp := &http.ServeMux{}
	muxApp := mux.NewRouter()
	pusher := eventsource.NewServer()

	muxApp.HandleFunc("/channel", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		bcastChan := pusher.Handler("hello_world")
		bcastChan.ServeHTTP(w, r)
	})

	server := &http.Server{
		Handler:      muxApp, // all this needs is something that fulfills the Handler interface
		Addr:         ":3001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	msgChan := make(chan eventPusherMessage, 1)

	// emit a message every two seconds
	go func() {
		for {
			time.Sleep(2 * time.Second)
			msgChan <- eventPusherMessage{
				EventStr: "message",
				Channel:  "hello_world",
				DataStr:  "{\"test\":\"message\"}",
			}
		}
	}()

	for {
		nextMsg := <-msgChan
		fmt.Println("go channel works")
		pusher.Publish([]string{"hello_world"}, &nextMsg)
	}

	server.ListenAndServe()

}
