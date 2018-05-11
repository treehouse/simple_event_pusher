package main

import (
	"fmt"
	es "github.com/donovanhide/eventsource"
	"net"
	"net/http"
	"time"
)

var incomingPort = ":8080"
var channel = "hello_world"

func main() {
	openTcpConn(func(tcpConn net.Listener) {
		pts := NewPushToSessionPipe()

		// assign new connections to uuid channel
		http.HandleFunc("/channel", pts.ServeSession(channel))

		// receive incoming cc results from redis
		// redis is contained here, nowhere else
		go generateMsgs(pts)

		// redis msgs push to browser
		// go pts.pushMsgs()

		// listen for new connections
		http.Serve(
			tcpConn,
			nil, /* Handler (DefaultServeMux if nil, could drop gorrilla in here) */
		)
	})
}

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
	pusher := es.NewServer()
	eventChannel := make(chan *eventPusherMessage, 1)
	return &PushToSession{
		eventPusher:     pusher,
		redisToPushChan: eventChannel,
		close:           pusher.Close,
	}
}

func (pts *PushToSession) ServeSession(sessionChannel string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer pts.close()
		setHeaders(&w)
		pushConn := pts.eventPusher.Handler(sessionChannel)
		// redis msgs push to browser
		go pts.pushMsgs()
		pushConn.ServeHTTP(w, r)
	}
}

func (pts *PushToSession) SendToPush(msg *eventPusherMessage) {
	pts.redisToPushChan <- msg
}

func (pts *PushToSession) pushMsgs() {
	for {
		nextMsg := <-pts.redisToPushChan
		fmt.Println("go channel works")
		pts.eventPusher.Publish([]string{nextMsg.Channel}, nextMsg)
	}
}

func openTcpConn(mux func(net.Listener)) {
	tcp, err := net.Listen("tcp", incomingPort)
	if err != nil {
		return
	}
	defer tcp.Close()

	// feed open tcp connection to potential multiplex
	mux(tcp)
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
