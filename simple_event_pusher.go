package main

import (
	"fmt"
	es "github.com/donovanhide/eventsource"
	"net"
	"net/http"
	"time"
	"regexp"
	"sync"
)

var incomingPort = ":8080"
var matchChannel = regexp.MustCompile(`^/channel/([^/]+)$`)

// connections are added to this map so redis can push to these connections
type Connections struct {
	list map[string]*PushToSession
	mu sync.Mutex
}

func (c *Connections) Add(sessionName string, pts *PushToSession) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.list[sessionName] = pts
}

func (c *Connections) Delete(sessionName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.list, sessionName)
}

func (c *Connections) SendToPush(channel string, msg *eventPusherMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if pts, ok := c.list[channel]; ok {
    pts.SendToPush(msg)
	}
}

func main() {
	openTcpConn(func(tcpConn net.Listener) {

		var connections = &Connections{list: map[string]*PushToSession{}}
		// assign new connections to uuid channel
		// need some type of second matcher here to split this url into two routes
		// a main channel route and a session subroute
		http.HandleFunc("/channel/", ServeSession(connections))

		// receive incoming cc results from redis
		// redis is contained here, nowhere else
		go generateMsgs(connections)

		// listen for new connections
		http.Serve(
			tcpConn,
			nil, /* Handler (DefaultServeMux if nil, gorrilla caused header problems) */
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

func NewPushToSessionPipe() *PushToSession {  // rename to New in separate package later
	pusher := es.NewServer()
	eventChannel := make(chan *eventPusherMessage, 1)
	return &PushToSession{
		eventPusher:     pusher,
		redisToPushChan: eventChannel,
		close:           pusher.Close,
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



func getChannel(r *http.Request) string {
	sessionChannelMatches := matchChannel.FindStringSubmatch(r.RequestURI)
	sessionChannel := sessionChannelMatches[1]
	fmt.Println(sessionChannel)
	return sessionChannel
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
}

func ServeSession(conns *Connections) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		
		
		sessionChannel := getChannel(r)
		setHeaders(&w)
		
		pts := NewPushToSessionPipe() // create pts conn for session
		defer pts.close()
		conns.Add(sessionChannel, pts)
		
		pushConn := pts.eventPusher.Handler(sessionChannel)
		// removes pts from session collection when connection is closed by browser
		defer conns.Delete(sessionChannel)

		// redis msgs push to browser
		go pts.pushMsgs()

		// opens event connection. 
		// When browser disconnects, eventPusher is closed
		// and function ends, ending pushMsgs goroutine
		pushConn.ServeHTTP(w, r) // blocking line 

		// runs until browser disconnects, then pts is cleaned up
	}
}

func openTcpConn(mux func(net.Listener)) {
	tcp, err := net.Listen("tcp", incomingPort)
	if err != nil {
		return
	}
	defer tcp.Close()

	mux(tcp)
}

func generateMsgs(conns *Connections) {
	binarySwitch := true // testing how this handles multiple channels
	for {
		time.Sleep(2 * time.Second)
		var channel string;
		if binarySwitch {
			channel = "hello-world"
		} else {
			channel = "hello_multiplex"
		}
		binarySwitch = !binarySwitch
		conns.SendToPush(channel, &eventPusherMessage{
			EventStr: "message",
			Channel:  channel,
			DataStr:  "{\"test\":\"message\"}",
		})
	}
}
