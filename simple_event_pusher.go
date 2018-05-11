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

type PushToSession struct {
	eventPusher     *es.Server
	redisToPushChan chan *eventPusherMessage
	close           func()
}

func NewPushToSessionPipe() *PushToSession {
	pusher, close := makeEvtPusher()
	eventChannel := make(chan *eventPusherMessage, 1)
	return &PushToSession{
		eventPusher:     pusher,
		redisToPushChan: eventChannel,
		close:           close,
	}
}

func (pts *PushToSession) ServeSession(sessionChannel string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer pts.close()
		setHeaders(&w)
		pushConn := pts.eventPusher.Handler(sessionChannel)
		pushConn.ServeHTTP(w, r)
	}
}

func (pts *PushToSession) SendToPush(msg *eventPusherMessage) {
	pts.redisToPushChan <- msg
}

func main() {
	openTcpConn(func(tcpConn net.Listener) {
		pts := NewPushToSessionPipe()

		// assign new connections to uuid channel
		http.HandleFunc("/channel", pts.ServeSession(channel))

		// listen for new connections
		go http.Serve(tcpConn, nil)

		// receive incoming cc results from redis
		go generateMsgs(pts)

		// redis msgs push to browser
		pushMsgs(pts)
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

func generateMsgs(pts *PushToSession) {
	for {
		time.Sleep(2 * time.Second)
		pts.SendToPush(&eventPusherMessage{
			EventStr: "message",
			Channel:  channel,
			DataStr:  "{\"test\":\"message\"}",
		})
	}
}

func pushMsgs(pts *PushToSession) {
	for {
		nextMsg := <-pts.redisToPushChan
		fmt.Println("go channel works")
		pts.eventPusher.Publish([]string{nextMsg.Channel}, nextMsg)
	}
}
