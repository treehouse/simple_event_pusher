package main

import (
	"net"
	"net/http"
	"time"
	mux "github.com/nicolasjhampton/simple_event_pusher/pkg/mux"
	tcp "github.com/nicolasjhampton/simple_event_pusher/pkg/tcp"
	event "github.com/nicolasjhampton/simple_event_pusher/pkg/event"
)


func main() {
	tcp.OpenPort(":8080", func(tcpConn net.Listener) {

		var connPusher = mux.New()
		// assign new connections to uuid channel
		// need some type of second matcher here to split this url into two routes
		// a main channel route and a session subroute
		http.HandleFunc("/channel/", connPusher.ServeSession)

		// receive incoming cc results from redis
		// redis is contained here, nowhere else
		go generateMsgs(connPusher)

		// listen for new connections
		http.Serve(
			tcpConn,
			nil, /* Handler (DefaultServeMux if nil, gorrilla caused header problems) */
		)
	})
}


//
func generateMsgs(s *mux.Server) {
	binarySwitch := true // testing how this handles multiple channels
	for {
		time.Sleep(2 * time.Second)
		var channel string;
		if binarySwitch {
			channel = "hello-world"
		} else {
			channel = "hello-multiplex"
		}
		binarySwitch = !binarySwitch
		SendToPush chan 
		cs.SendToPush <- SendStruct{
			Key: channel, 
			Msg: &event.Message{
				EventStr: "message",
				Channel:  channel,
				DataStr:  "{\"test\":\"message\"}",
			}
		}
	}
}