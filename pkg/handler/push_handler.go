package handler

import (
	push "github.com/treehouse/simple_event_pusher/internal/connection"
	mux "github.com/treehouse/simple_event_pusher/pkg/push_mux"
	"net/http"
	"regexp"
)

var channelRegex = regexp.MustCompile(`^/channel/([^/]+)$`)

func getChannel(r *http.Request) string {
	return channelRegex.FindStringSubmatch(r.RequestURI)[1]
}

func ServeSession(cs *mux.ConnStore, cors string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if cors != "" {
			w.Header().Add("Access-Control-Allow-Origin", cors)
		}

		sessionChannel := getChannel(r)

		conn := push.NewConnection()
		defer conn.Close()

		cs.Add(sessionChannel, conn)
		defer cs.Delete(sessionChannel)

		push := conn.Handler(sessionChannel)

		// receives redis msgs and pushes to browser, one thread per session
		go conn.Msgs()

		// opens event connection. blocking line
		push.ServeHTTP(w, r)
		// when browser disconnects, connection is closed and
		// removed from ConnStore, and goroutine for that browser's
		// push event msgs ends
	}
}
