package push

import (
	"fmt"
	es "github.com/donovanhide/eventsource"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"net/http"
)

type Connection struct {
	eventPusher     *es.Server
	redisToPushChan chan *event.Message
	close           func()
}

func NewConnection() *Connection { // rename to New in separate package later
	pusher := es.NewServer()
	eventChannel := make(chan *event.Message, 1)
	return &Connection{
		eventPusher:     pusher,
		redisToPushChan: eventChannel,
		close:           pusher.Close,
	}
}

func (c *Connection) Handler(sessionChannel string) http.HandlerFunc {
	return c.eventPusher.Handler(sessionChannel)
}

func (c *Connection) Close() {
	c.close()
}

func (c *Connection) SendToPush(msg *event.Message) {
	c.redisToPushChan <- msg
}

func (c *Connection) Msgs() {
	for {
		nextMsg := <-c.redisToPushChan
		fmt.Println(nextMsg)
		c.eventPusher.Publish([]string{nextMsg.Channel}, nextMsg)
	}
}
