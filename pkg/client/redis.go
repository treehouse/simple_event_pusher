package client

import (
	event "github.com/treehouse/simple_event_pusher/internal/event"
	mux "github.com/treehouse/simple_event_pusher/pkg/push_mux"
	"github.com/go-redis/redis"
	"encoding/json"
	"fmt"
)

type Incoming interface {
	Close() error
	ReceiveMessage() (*redis.Message, error)
}

func Redis(redisAddr string, redisChannel string) Incoming {
	rClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pubsub := rClient.Subscribe(redisChannel)
	return pubsub
}

// PUBLISH event_pusher '{ "event": "", "channel": "hello-multiplex", "data": "this is data"}'
func ListenForMsgs(cs *mux.ConnStore, dataSource Incoming) {
	
	defer dataSource.Close()

	for {
		msg, err := dataSource.ReceiveMessage();
		epmsg := event.Message{} 
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("got payload from redis pubsub %+v\n", msg.Payload)
		if err := json.Unmarshal([]byte(msg.Payload), &epmsg); err != nil {
			fmt.Println(err)
			continue
		}

		go cs.Send(&epmsg)
	}
}