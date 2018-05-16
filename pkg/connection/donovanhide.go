package push

import (
	es "github.com/donovanhide/eventsource"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"net/http"
	"fmt"
)

// Manages the one-way push connection to a single browser. Enevlops
// and isolates push method, currently using donovanhide/eventsource
// library under the hood, which runs perfectly but spawns a multiplex
// push server for each single user channel that we can't monitor or log: 
// https://godoc.org/github.com/donovanhide/eventsource
//
// TODO: Consider adding a push method with ie/edge support:
//
// https://caniuse.com/#search=eventsource
//
// https://caniuse.com/#search=websockets
type DonovanhideConnection struct {
	channel     string
	handler http.HandlerFunc
	toPushChan  chan event.Message
	// Closed bool
	eventPusher *es.Server
}

func NewDonovanConn(sessionChannel string) Connection {
	pusher := es.NewServer()
	handler := pusher.Handler(sessionChannel) // donovanhide itself is a mux
	eventChannel := make(chan event.Message)
	return &DonovanhideConnection{
		channel:     sessionChannel,
		handler: handler,
		eventPusher: pusher,
		toPushChan:  eventChannel,
	}
}

// ServePUSH is a wrapper around http.HandlerFunc.ServeHTTP. Starts a
// keep-alive connection channel to a single browser connection.
// ServePUSH will block execution in the thread its in until browser
// disconnects.
func (c *DonovanhideConnection) ServePUSH(w http.ResponseWriter, r *http.Request) {
	
	c.handler.ServeHTTP(w, r)
}

// Send provides a new message to the running Msgs
// thread to push to the browser.
func (c *DonovanhideConnection) Send(msg event.Message) {
	fmt.Println(msg)
	c.toPushChan <- msg
}

// Returns the channel this connection is assigned to in the connList
func (c *DonovanhideConnection) Channel() string {
	return c.channel
}

// Responsible for cleaning up the connection after disconnect. Should 
// be called with defer when Msgs is used to clean up disconnectioned 
// broswer connections.
// TODO: Improve to match CustomConnection as closely as possible
func (c *DonovanhideConnection) Close() {
	close(c.toPushChan)
}


// Msgs is meant to be ran on it's own goroutine and listen
// for any messages to push to the browser. Run with a deferred
// Close function to clean up disconnected browser connections.
// TODO: Check to see if this method is nessasary.
func (c *DonovanhideConnection) Msgs() {
	// notice that this closes the multiplex server running for one our
	// one connection, not just the connection
	defer c.eventPusher.Close() 
	for {
		nextMsg, ok := <-c.toPushChan
		if !ok {
			return
		}
		fmt.Println("pushing to browser:", nextMsg)

		c.eventPusher.Publish([]string{nextMsg.GetChannel()}, nextMsg)
	}
}