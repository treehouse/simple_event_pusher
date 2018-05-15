package push

import (
	"fmt"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	es "github.com/donovanhide/eventsource"
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
type CustomConnection struct {
	channel     string
	handler http.HandlerFunc
	toPushChan  chan event.Message
	Closed bool
	// eventPusher *es.Server
}

func NewCustomConn(sessionChannel string) Connection {
	newConn := &CustomConnection{
		channel:     sessionChannel,
		Closed: false,
		// handler: handler,
		// toPushChan:  eventChannel,
		// eventPusher: pusher,
	}
	// pusher := es.NewServer()
	// handler := pusher.Handler(sessionChannel)
	newConn.handler = newConn.NewHandler(/* cors string */)
	newConn.toPushChan = make(chan event.Message)
	return newConn
	// return &CustomConnection{
	// 	channel:     sessionChannel,
	// 	handler: handler,
	// 	toPushChan:  eventChannel,
	// 	// eventPusher: pusher,
	// }
}

func (c *CustomConnection) NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		h := w.Header()
		h.Set("Content-Type", "text/event-stream; charset=utf-8")
		h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		h.Set("Connection", "keep-alive")

		// Incorporate with the Access-Control-Allow-Origin env variable
		h.Set("Access-Control-Allow-Origin", "*")

		// holding off on gzip for now
		// useGzip := srv.Gzip && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
		// if useGzip {
		// 	h.Set("Content-Encoding", "gzip")
		// }

		w.WriteHeader(http.StatusOK)

		// If the Handler is still active even though the connection is closed, stop here.
		// Otherwise the Handler may block indefinitely.
		if c.Closed == true {
			return
		}

		// sub := &subscription{
		// 	channel:     channel,
		// 	lastEventId: req.Header.Get("Last-Event-ID"),
		// 	out:         make(chan Event, srv.BufferSize),
		// }


		// srv.subs <- sub

		flusher := w.(http.Flusher)
		notifier := w.(http.CloseNotifier)

		flusher.Flush()
		
		// Still using donovanhide here for the encoder
		// holding off on gzip for now
		useGzip := false
		enc := es.NewEncoder(w, useGzip)


		for {
			select {
			case <-notifier.CloseNotify():
				// accomodate for multiple users on a channel here potentially
				// not currently a concern
				c.Close()
				return
			case ev, ok := <-c.toPushChan:
				if !ok {
					return
				}
				if err := enc.Encode(ev); err != nil {
					// accomodate for multiple users on a channel here potentially
					// not currently a concern
					c.Close()
					// srv.unregister <- sub

					fmt.Println(err)
					// TODO: put logger back on connection struct
					// if srv.Logger != nil {
					// 	srv.Logger.Println(err)
					// }
					return
				}
				// Flusher flushes the encoder's (enc) buffer to the client
				// it's a http library interface
				flusher.Flush()
			}
		}
	}
}

// ServePUSH is a wrapper around http.HandlerFunc.ServeHTTP. Starts a
// keep-alive connection channel to a single browser connection.
// ServePUSH will block execution in the thread its in until browser
// disconnects.
func (c *CustomConnection) ServePUSH(w http.ResponseWriter, r *http.Request) {
	
	c.handler.ServeHTTP(w, r)
}

// Send provides a new message to the running Msgs
// thread to push to the browser.
func (c *CustomConnection) Send(msg event.Message) {
	fmt.Println(msg)
	c.toPushChan <- msg
}

// Msgs is meant to be ran on it's own goroutine and listen
// for any messages to push to the browser. Run with a deferred
// Close function to clean up disconnected browser connections.
func (c *CustomConnection) Msgs() {
	// defer c.Close()
	// for {
	// 	nextMsg, ok := <-c.toPushChan
	// 	if !ok {
	// 		return
	// 	}
	// 	fmt.Println("pushing to browser:", nextMsg)

	// 	c.handler.Publish([]string{nextMsg.GetChannel()}, nextMsg)
	// }
}

// Close is a wrapper around the implementation of push. Responsible
// for cleaning up the connection after disconnect. Should be called
// with defer when Msgs is used to clean up disconnectioned broswer
// connections. Setting toPushChan to nil so handler can read nil
// value and know the connection is no longer open to handle.
func (c *CustomConnection) Close() {
	close(c.toPushChan)
	c.Closed = true
}

// Returns the channel this connection is assigned to in the connList
func (c *CustomConnection) Channel() string {
	return c.channel
}
