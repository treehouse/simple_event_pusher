package main

import (
	event "github.com/nicolasjhampton/simple_event_pusher/pkg/event"
	handler "github.com/nicolasjhampton/simple_event_pusher/pkg/handler"
	mux "github.com/nicolasjhampton/simple_event_pusher/pkg/push_mux"
	tcp "github.com/nicolasjhampton/simple_event_pusher/pkg/tcp"
	"net"
	"net/http"
	"time"
)

func main() {
	tcp.OpenPort(":8080", func(tcpConn net.Listener) {

		var connMux = mux.New()

		// runs on a single thread and manages connection list
		// access as many threads read/write connections via channels
		// to this thread
		go connMux.Run()
		// assign new connections to uuid channel
		// need some type of second matcher here to split this url into two routes
		// a main channel route and a session subroute
		http.HandleFunc("/channel/", handler.ServeSession(connMux))

		// receive incoming cc results from redis
		// redis is contained here, nowhere else
		go generateMsgs(connMux)

		// listen for new connections
		http.Serve(
			tcpConn,
			nil, /* Handler (DefaultServeMux if nil) */
		)
	})
}

//
func generateMsgs(cs *mux.ConnStore) {
	binarySwitch := true // testing how this handles multiple channels
	for {
		/** mock msgs code, replace with redis */
		time.Sleep(2 * time.Second)
		var channel string
		if binarySwitch {
			channel = "hello-world"
		} else {
			channel = "hello-multiplex"
		}
		binarySwitch = !binarySwitch
		/***************************************/

		cs.Send(channel, &event.Message{
			EventStr: "message",
			Channel:  channel,
			DataStr:  "{\"test\":\"message\"}",
		})
	}
}
