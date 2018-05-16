package push

import (
	"fmt"
	es "github.com/donovanhide/eventsource"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"net/http"
)

// Manages the one-way push connection to a single browser. Enevlops
// and isolates push method, currently using minimal sections of the
// donovanhide/eventsource library for encoding event bodies under
// the hood: https://godoc.org/github.com/donovanhide/eventsource
//
// TODO: Consider adding a push method with ie/edge support:
//
// https://caniuse.com/#search=eventsource
//
// https://caniuse.com/#search=websockets
type CustomConnection struct {
	channel    string
	handler    http.HandlerFunc
	toPushChan chan event.Message
	Closed     bool
}

func NewCustomConn(sessionChannel, cors string) Connection {
	newConn := &CustomConnection{
		channel: sessionChannel,
		Closed:  false,
	}
	newConn.handler = newConn.NewHandler(cors)
	newConn.toPushChan = make(chan event.Message)
	return newConn
}

func (c *CustomConnection) NewHandler(cors string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		h := w.Header()
		h.Set("Content-Type", "text/event-stream; charset=utf-8")
		h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		h.Set("Connection", "keep-alive")
		h.Set("Access-Control-Allow-Origin", cors)

		// holding off on gzip for now
		// useGzip := srv.Gzip && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
		// if useGzip {
		// 	h.Set("Content-Encoding", "gzip")
		// }

		w.WriteHeader(http.StatusOK)

		// If the Handler is still active even though the connection is
		// closed, stop here. Otherwise the Handler may block indefinitely.
		if c.Closed == true {
			return
		}

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
				// Convert to logger, maybe still use fmt for doctests
				fmt.Println("client disconnected from channel: ", c.Channel())
				// Accomodate for multiple users on a channel here potentially.
				// Not currently a concern.
				c.Close()
				return
			case ev, ok := <-c.toPushChan:
				if !ok {
					// Might need a c.Close call here. Investigating.
					return
				}
				if err := enc.Encode(ev); err != nil {
					// Accomodate for multiple users on a channel here potentially.
					// Not currently a concern.
					c.Close()
					fmt.Println(err)
					// TODO: put logger back on connection struct
					// if srv.Logger != nil {
					// 	srv.Logger.Println(err)
					// }
					return
				}
				// Flusher flushes the encoder's (enc) buffer to the client.
				// It's an http library interface.
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

// Send provides a new message to the for loop running in our custom
// http.HandlerFunc (see NewHandler).
func (c *CustomConnection) Send(msg event.Message) {
	fmt.Println(msg)
	c.toPushChan <- msg
}

// Returns the channel this connection is assigned to in the connList
func (c *CustomConnection) Channel() string {
	return c.channel
}

// Responsible for cleaning up the connection after disconnect.
// CustomConnection calls Close automatically when a browser
// disconnects to clean up any resources. By setting Closed to
// true, we avoid other browsers arriving at this route while it's
// closing with a check in our HandlerFunc.
func (c *CustomConnection) Close() {
	close(c.toPushChan)
	c.Closed = true
}

// So far, msgs is not needed by CustomConnection, but is required
// by the Connection interface.
// TODO: Consider factoring out the for loop in http.HandlerFunc
// to use it for easier testing.
func (c *CustomConnection) Msgs() {

}
