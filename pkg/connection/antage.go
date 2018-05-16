package push

import (
	es "gopkg.in/antage/eventsource.v1"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"net/http"
	"fmt"
)

// Manages the one-way push connection to a single browser. Enevlops
// and isolates push method. Antage library is lightweight and does 
// not create a mutliplex of it's own, but current implementation 
// creating a dependable connection: https://github.com/antage/eventsource
//
// TODO: Consider adding a push method with ie/edge support:
//
// https://caniuse.com/#search=eventsource
//
// https://caniuse.com/#search=websockets
type AntageConnection struct {
	channel     string
	handler es.EventSource
	toPushChan  chan event.Message
	// Closed bool
}

func NewAntageConn(sessionChannel string) Connection {
	handler := es.New(nil, nil)
	eventChannel := make(chan event.Message, 1)
	return &AntageConnection{
		channel:     sessionChannel,
		handler: handler,
		toPushChan:  eventChannel,
	}
}

// ServePUSH is a wrapper around http.HandlerFunc.ServeHTTP. Starts a
// keep-alive connection channel to a single browser connection.
// ServePUSH will block execution in the thread its in until browser
// disconnects.
func (c *AntageConnection) ServePUSH(w http.ResponseWriter, r *http.Request) {
	c.handler.ServeHTTP(w, r)
}

// Send provides a new message to the running Msgs thread to push to
// the browser.
func (c *AntageConnection) Send(msg event.Message) {
	c.toPushChan <- msg
}
// Returns the channel this connection is assigned to in the connList
func (c *AntageConnection) Channel() string {
	return c.channel
}

// Responsible for cleaning up the connection after disconnect. 
// Should be called with defer when Msgs is used to clean up 
// disconnectioned broswer connections.
// TODO: Improve to match CustomConnection as closely as possible
func (c *AntageConnection) Close() {
	// One disconnecting browser closes all browser connections on this
	// channel. Fine for now as only one user is on a channel, but this
	// package is more thqan capable of managing private chat rooms as
	// channels if we count users currently left on this channel
	close(c.toPushChan)
}

// Msgs is meant to be ran on it's own goroutine and listen for any
// messages to push to the browser.
func (c *AntageConnection) Msgs() {
	for {
		nextMsg, ok := <-c.toPushChan
		fmt.Println(ok) // <- this channel is closing at start for some reason
		if !ok {
			// when the channel is closed, for loop is broken, eventPusher
			// is closed, and goroutine ends.
			// TODO: may want to call c.Close here instead of in the handler.
			return
		}
		fmt.Println("pushing to browser:", nextMsg)
		c.handler.SendEventMessage(nextMsg.Data(), nextMsg.Event(), nextMsg.Id())
	}
}
