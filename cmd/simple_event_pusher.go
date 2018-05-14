package main

import (
	cli "github.com/treehouse/simple_event_pusher/pkg/client"
	push "github.com/treehouse/simple_event_pusher/pkg/connection"
	env "github.com/treehouse/simple_event_pusher/pkg/env"
	handler "github.com/treehouse/simple_event_pusher/pkg/handler"
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
	redisPubsubChannel := env.Default("EVENT_PUSHER_REDIS_PUBSUB_CHANNEL", DEFAULT_REDIS_PUBSUB_CHANNEL)
	accessControlAllowOrigin := env.Default("EVENT_PUSHER_ACCESS_CONTROL_ALLOW_ORIGIN", DEFAULT_ACCESS_CONTROL_ALLOW_ORIGIN)

	s := &http.Server{
		Addr:        port,
		ReadTimeout: 15 * time.Second,
	}

	// Threadsafe collection of connections
	connList := push.NewConnList()

	// creates a new thread for each new session connection
	http.HandleFunc("/channel/", handler.ServeSession(connList, accessControlAllowOrigin))

	// all redis msgs come in on a single thread
	go cli.ListenForMsgs(
		connList,
		cli.Redis(redisAddr, redisPubsubChannel),
	)

	// Listen for new connections
	// https://golang.org/pkg/net/http/#Server.ListenAndServe
	// https://golang.org/pkg/net/http/#Server.Serve
	s.ListenAndServe()
}
