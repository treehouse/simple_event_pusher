package mux

import (
	"regexp"
	"sync"
	"net/http"
	push "github.com/nicolasjhampton/simple_event_pusher/internal/connection"
	event "github.com/nicolasjhampton/simple_event_pusher/pkg/event"
)

var channelRegex = regexp.MustCompile(`^/channel/([^/]+)$`)


type connList map[string]*push.Connection

type Server struct {
	list connList
	mu sync.RWMutex
}

func New() *Server {
	return &Server{
		list: connList{},
	}
}

func (s *Server) ServeSession(w http.ResponseWriter, r *http.Request) {
	sessionChannel := getChannel(r)
	setHeaders(&w)
	
	conn := push.NewConnection() // create pts conn for session
	defer conn.Close()
	s.Add(sessionChannel, conn)
	
	push := conn.Handler(sessionChannel)
	// removes pts from session collection when connection is closed by browser
	defer s.Delete(sessionChannel)

	// redis msgs push to browser
	go conn.Msgs()

	// opens event connection. 
	// When browser disconnects, eventPusher is closed
	// and function ends, ending pushMsgs goroutine
	push.ServeHTTP(w, r) // blocking line 

	// runs until browser disconnects, then pts is cleaned up
}

func (s *Server) Add(sessionName string, conn *push.Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.list[sessionName] = conn
}

func (s *Server) Delete(sessionName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.list, sessionName)
}

func (s *Server) SendToPush(sessionName string, msg *event.Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if conn, ok := s.list[sessionName]; ok {
    conn.SendToPush(msg)
	}
}

func getChannel(r *http.Request) string {
	return channelRegex.FindStringSubmatch(r.RequestURI)[1]
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
}