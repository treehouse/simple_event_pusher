package main

import (
	"fmt"
	es "github.com/donovanhide/eventsource"
	// "github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

var incomingPort = ":8080";
var channel = "hello_world"

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
	openTcpConn(func(tcpConn net.Listener) {
		eventPusher, close := makeEvtPusher();
		msgChan := make(chan eventPusherMessage, 1)

		// assign new connections to uuid channel
		http.HandleFunc("/channel", makeHandlerFunc(eventPusher, close))

		// listen for new connections
		go http.Serve(tcpConn, nil)

		// receive incoming cc results from redis
		go generateMsgs(msgChan)

		// redis msgs push to browser
		pushMsgs(msgChan, eventPusher)
	})
}

func openTcpConn(f func(net.Listener)) {
	tcp, err := net.Listen("tcp", incomingPort)
	if err != nil {
		return
	}
	defer tcp.Close()
	f(tcp)
}



func makeEvtPusher() (*es.Server, func()) {
	pusher := es.NewServer()
	return pusher, pusher.Close
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/json")
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
}

func makeHandlerFunc(pushServer *es.Server, close func()) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer close()
		setHeaders(&w)
		bcastChan := pushServer.Handler(channel)
		bcastChan.ServeHTTP(w, r)
	}
}

func generateMsgs(pipeToPush chan eventPusherMessage) {
	for {
		time.Sleep(2 * time.Second)
		pipeToPush <- eventPusherMessage{
			EventStr: "message",
			Channel:  channel,
			DataStr:  "{\"test\":\"message\"}",
		}
	}
}

func pushMsgs(msgs chan eventPusherMessage, eventer *es.Server) {
	for {
		nextMsg := <-msgs
		fmt.Println("go channel works")
		eventer.Publish([]string{nextMsg.Channel}, &nextMsg)
	}
}
