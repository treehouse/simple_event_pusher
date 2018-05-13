package main

import (
	cli "github.com/treehouse/simple_event_pusher/pkg/client"
	handler "github.com/treehouse/simple_event_pusher/pkg/handler"
	mux "github.com/treehouse/simple_event_pusher/pkg/push_mux"
	"net/http"
	"os"
	"time"
)

const (
	DEFAULT_PORT                        = ":8080"
	DEFAULT_REDIS_ADDR                  = "localhost:6379"
	DEFAULT_REDIS_PUBSUB_CHANNEL        = "event_pusher"
	DEFAULT_ACCESS_CONTROL_ALLOW_ORIGIN = "http://localhost:3001"
)

func main() {
	port := envDefault("EVENT_PUSHER_PORT", DEFAULT_PORT)
	redisAddr := envDefault("EVENT_PUSHER_REDIS_ADDR", DEFAULT_REDIS_ADDR)
	redisPubsubChannel := envDefault("EVENT_PUSHER_REDIS_PUBSUB_CHANNEL", DEFAULT_REDIS_PUBSUB_CHANNEL)
	accessControlAllowOrigin := envDefault("EVENT_PUSHER_ACCESS_CONTROL_ALLOW_ORIGIN", DEFAULT_ACCESS_CONTROL_ALLOW_ORIGIN)

	s := &http.Server{
		Addr:        port,
		ReadTimeout: 15 * time.Second,
	}

	// runs on a single thread, manages connection list
	// data structure uses a RWMutex to provide concurrent reads
	// and threadsafe on writes
	connMux := mux.New()
	go connMux.Run()

	// creates a new thread for each new session connection
	http.HandleFunc("/channel/", handler.ServeSession(connMux, accessControlAllowOrigin))

	// all redis msgs come in on a single thread
	// and leave on separate goroutines
	go cli.ListenForMsgs(
		connMux,
		cli.Redis(redisAddr, redisPubsubChannel),
	)

	s.ListenAndServe()
}

// PUBLISH event_pusher '{ "event": "", "channel": "hello-multiplex", "data": "this is data"}'

func envDefault(key, def string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return def
	}
	return val
}
