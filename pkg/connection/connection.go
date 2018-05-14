package push

import (
	"fmt"
	es "github.com/donovanhide/eventsource"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"net/http"
)

// Manages the one-way push connection to a single browser. Enevlops
// and isolates push method, currently using donovanhide/eventsource
// library under the hood: https://godoc.org/github.com/donovanhide/eventsource
//
// TODO: Consider switching to a push method with ie/edge support:
//
// https://caniuse.com/#search=eventsource
//
// https://caniuse.com/#search=websockets
type Connection struct {
	Channel     string
	pushHandler http.HandlerFunc
	eventPusher *es.Server
	toPushChan  chan *event.Message
	close       func()
}

func NewConnection(sessionChannel string) *Connection {
	pusher := es.NewServer()
	handler := pusher.Handler(sessionChannel)
	eventChannel := make(chan *event.Message, 1)
	return &Connection{
		Channel:     sessionChannel,
		pushHandler: handler,
		eventPusher: pusher,
		toPushChan:  eventChannel,
		close:       pusher.Close,
	}
}

// ServePUSH is a wrapper around http.HandlerFunc.ServeHTTP. Starts a
// keep-alive connection channel to a single browser connection.
// ServePUSH will block execution in the thread its in until browser
// disconnects.
func (c *Connection) ServePUSH(w http.ResponseWriter, r *http.Request) {
	c.pushHandler.ServeHTTP(w, r)
}

// SendToPush provides a new message to the running Msgs
// thread to push to the browser.
func (c *Connection) Send(msg *event.Message) {
	c.toPushChan <- msg
}

// Msgs is meant to be ran on it's own goroutine and listen
// for any messages to push to the browser. Run with a deferred
// Close function to clean up disconnected browser connections.
func (c *Connection) Msgs() {
	for {
		nextMsg := <-c.toPushChan
		fmt.Println("pushing to browser:", nextMsg)
		c.eventPusher.Publish([]string{nextMsg.Channel}, nextMsg)
	}
}

// Close is a wrapper around the implementation of push. Responsible
// for cleaning up the connection after disconnect. Should be called
// with defer when Msgs is used to clean up disconnectioned broswer
// connections.
func (c *Connection) Close() {
	c.close()
}
