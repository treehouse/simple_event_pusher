package handler

import (
	"net/http"
	"regexp"
	push "github.com/nicolasjhampton/simple_event_pusher/internal/connection"
	mux "github.com/nicolasjhampton/simple_event_pusher/pkg/push_mux"
)

var channelRegex = regexp.MustCompile(`^/channel/([^/]+)$`)

func getChannel(r *http.Request) string {
	return channelRegex.FindStringSubmatch(r.RequestURI)[1]
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
}

func ServeSession(cs *mux.ConnStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionChannel := getChannel(r)
		setHeaders(&w)
		
		conn := push.NewConnection(); defer conn.Close();

		cs.Add(sessionChannel, conn); defer cs.Delete(sessionChannel);
		
		push := conn.Handler(sessionChannel)
		// redis msgs push to browser
		go conn.Msgs()

		// opens event connection. 
		push.ServeHTTP(w, r) // blocking line 
		// when browser disconnects, conn, cs, and goroutine are cleaned up
	}
}