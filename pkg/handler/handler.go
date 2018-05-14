package handler

import (
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	// mux "github.com/treehouse/simple_event_pusher/pkg/push_mux"
	"net/http"
	"regexp"
)

var chanRegex = regexp.MustCompile(`^/channel/([^/]+)$`)

func getChannel(r *http.Request) string {
	if channelSlice := chanRegex.FindStringSubmatch(r.RequestURI); len(channelSlice) > 1 {
		return channelSlice[1]
	}
	return ""
}

func ServeSession(cl *push.ConnList, cors string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionChannel := getChannel(r)
		if sessionChannel == "" {
			return
		}

		if cors != "" {
			w.Header().Add("Access-Control-Allow-Origin", cors)
		}

		pConn := push.NewConnection(sessionChannel)
		defer pConn.Close()

		cl.Add(pConn)
		defer cl.Remove(pConn)

		go pConn.Msgs()

		pConn.ServePUSH(w, r)
	}
}
