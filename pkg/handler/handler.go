package handler

import (
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	mux "github.com/treehouse/simple_event_pusher/pkg/mux"
	"net/http"
	"regexp"
	"fmt"
)

var chanRegex = regexp.MustCompile(`^/v1/channels/([^/]+)$`)

// Private method that reads the channel from the url. Used by
// ServeSession.
func getChannel(r *http.Request) string {
	if channelSlice := chanRegex.FindStringSubmatch(r.RequestURI); len(channelSlice) > 1 {
		return channelSlice[1]
	}
	return ""
}

// ServeSession returns a function that can be used as a http.Handler
// to the http.DefaultServeMux. The http.Server in the http
// library will spawn a new goroutine for each request for a new
// keep-alive connection. http.Server will call the http.DefaultServeMux,
// which will call ServeSession to get a handler function for the
// specific requested session channel name in the url. This handler
// creates a new Connection for the session channel, spawns a goroutine
// to listen for messages for that channel from the redis client,
// and opens the keep-alive event source connection to the browser
// with ServePUSH.
func ServeSession(cs *mux.ConnStore, cors string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionChannel := getChannel(r)
		if sessionChannel == "" {
			return
		}

		// if cors != "" {
		// 	w.Header().Add("Access-Control-Allow-Origin", cors)
		// }

		pConn := push.NewCustomConn(sessionChannel, cors)
		//defer pConn.Close()

		cs.Add(pConn)
		defer cs.Remove(pConn)

		go pConn.Msgs()

		pConn.ServePUSH(w, r)
		fmt.Println("passed where I should")
	}
}
