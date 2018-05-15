package push

import (
	"fmt"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	es "gopkg.in/antage/eventsource.v1"
	"net/http"
)

// type Connection interface {
// 	ServePUSH(http.ResponseWriter, *http.Request)
// 	Send(*event.Message)
// 	Channel() string
// 	Close()
// 	Msgs()
// }

// Manages the one-way push connection to a single browser. Enevlops
// and isolates push method, currently using donovanhide/eventsource
// library under the hood: https://godoc.org/github.com/donovanhide/eventsource
//
// TODO: Consider switching to a push method with ie/edge support:
//
// https://caniuse.com/#search=eventsource
//
// https://caniuse.com/#search=websockets
//
// FIXED: Switched to https://github.com/antage/eventsource as an
// underlying eventsource implementation, which allows us to implement
// our own multiplex.
// TODO: Decouple any underlying implementation with an interface.
type AntageConnection struct {
	channel     string
	eventPusher es.EventSource
	toPushChan  chan event.Message
}

func NewAntageConn(sessionChannel string) Connection {
	pusher := es.New(nil, nil)
	eventChannel := make(chan event.Message, 1)
	return &AntageConnection{
		channel:     sessionChannel,
		eventPusher: pusher,
		toPushChan:  eventChannel,
	}
}

// Msgs is meant to be ran on it's own goroutine and listen for any
// messages to push to the browser. Run with a deferred Close function
// to clean up disconnected browser connections.
func (c *AntageConnection) Msgs() {
	defer c.eventPusher.Close()
	for {
		nextMsg, ok := <-c.toPushChan
		if !ok {
			// when the channel is closed, for loop is broken, eventPusher
			// is closed, and goroutine ends.
			return
		}
		fmt.Println("pushing to browser:", nextMsg)
		// "data-payload"
		c.eventPusher.SendEventMessage(nextMsg.Data(), nextMsg.Event(), nextMsg.Id())
	}
}

// ServePUSH is a wrapper around http.HandlerFunc.ServeHTTP. Starts a
// keep-alive connection channel to a single browser connection.
// ServePUSH will block execution in the thread its in until browser
// disconnects.
func (c *AntageConnection) ServePUSH(w http.ResponseWriter, r *http.Request) {
	c.eventPusher.ServeHTTP(w, r)
}

// Close is a wrapper around the implementation of push. Responsible
// for cleaning up the connection after disconnect. Should be called
// with defer when Msgs is used to clean up disconnectioned broswer
// connections.
func (c *AntageConnection) Close() {
	// One disconnecting browser closes all browser connections on this
	// channel. Fine for now as only one user is on a channel, but this
	// package is more thqan capable of managing private chat rooms as
	// channels if we count users currently left on this channel
	close(c.toPushChan)
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
