package push

import (
	"fmt"
	es "github.com/donovanhide/eventsource"
	event "github.com/treehouse/simple_event_pusher/internal/event"
	"net/http"
)

// Manages the one-way push connection to a single browser. Enevlops
// and isolates push method, currently using donovanhide/eventsource
type Connection struct {
	eventPusher     *es.Server
	toPushChan chan *event.Message
	close           func()
}

func NewConnection() *Connection {
	pusher := es.NewServer()
	eventChannel := make(chan *event.Message, 1)
	return &Connection{
		eventPusher:     pusher,
		toPushChan: eventChannel,
		close:           pusher.Close,
	}
}

// Handler is a initialization wrapper around the 
// implementation of push. Initializes a channel name
// to a single browser connection.
func (c *Connection) Handler(sessionChannel string) http.HandlerFunc {
	return c.eventPusher.Handler(sessionChannel)
}

// SendToPush provides a new message to the running Msgs
// thread to push to the browser.
func (c *Connection) SendToPush(msg *event.Message) {
	c.toPushChan <- msg
}

// Msgs is meant to be ran on it's own goroutine and listen
// for any messages to push to the browser
func (c *Connection) Msgs() {
	for {
		nextMsg := <-c.toPushChan
		fmt.Println("pushing to browser:", nextMsg)
		c.eventPusher.Publish([]string{nextMsg.Channel}, nextMsg)
	}
}

// Close is a wrapper around the implementation of push. Responsible
// for cleaning up the connection after disconnect.
func (c *Connection) Close() {
	c.close()
}
