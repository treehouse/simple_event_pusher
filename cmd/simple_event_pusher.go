package main

import (
	cli "github.com/treehouse/simple_event_pusher/pkg/client"
	env "github.com/treehouse/simple_event_pusher/pkg/env"
	handler "github.com/treehouse/simple_event_pusher/pkg/handler"
	mux "github.com/treehouse/simple_event_pusher/pkg/mux"
	"net/http"
	"time"
)

const (
	DEFAULT_PORT                        = ":8080"
	DEFAULT_REDIS_ADDR                  = "localhost:6379"
	DEFAULT_REDIS_PUBSUB_CHANNEL        = "event_pusher"
	DEFAULT_ACCESS_CONTROL_ALLOW_ORIGIN = "http://localhost:3001"
)

func main() {
	port := env.Default("EVENT_PUSHER_PORT", DEFAULT_PORT)
	redisAddr := env.Default("EVENT_PUSHER_REDIS_ADDR", DEFAULT_REDIS_ADDR)
	redisPubsubPrefix := envDefault("DEFAULT_REDIS_PUBSUB_CHANNEL_PREFIX", DEFAULT_REDIS_PUBSUB_CHANNEL_PREFIX)
	// redisPubsubChannel := env.Default("EVENT_PUSHER_REDIS_PUBSUB_CHANNEL", DEFAULT_REDIS_PUBSUB_CHANNEL)
	accessControlAllowOrigin := env.Default("EVENT_PUSHER_ACCESS_CONTROL_ALLOW_ORIGIN", DEFAULT_ACCESS_CONTROL_ALLOW_ORIGIN)

	s := &http.Server{
		Addr:        port,
		ReadTimeout: 15 * time.Second,
	}

	// Threadsafe collection of connections
	connStore := mux.New()

	// creates a new thread for each new session connection
	http.HandleFunc("/v1/channels/", handler.ServeSession(connStore, accessControlAllowOrigin))

	// Separate route for load balancer health check.
	http.HandleFunc("/v1/healthz", handleHealthCheck)

	// all redis msgs come in on a single thread
	go cli.ListenForMsgs(
		connStore,
		cli.Redis(redisAddr, redisPubsubPrefix),
	)

	// Listen for new connections
	// https://golang.org/pkg/net/http/#Server.ListenAndServe
	// https://golang.org/pkg/net/http/#Server.Serve
	s.ListenAndServe()
}

// Completely separate handler for load balancer health check.
// Not significant part of main data flow.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {

	// Write writes the data to the connection as part of an HTTP reply.
	//
	// If WriteHeader has not yet been called, Write calls
	// WriteHeader(http.StatusOK) before writing the data. If the Header
	// does not contain a Content-Type line, Write adds a Content-Type set
	// to the result of passing the initial 512 bytes of written data to
	// DetectContentType.

	// w.WriteHeader(200)
	// w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

