package client

import (
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	mux "github.com/treehouse/simple_event_pusher/pkg/push_mux"
	"github.com/go-redis/redis"
	"encoding/json"
	"fmt"
)

// PUBLISH event_pusher '{ "event": "", "channel": "hello-multiplex", "data": "this is data"}'
func ListenForMsgs(cs *mux.ConnStore, redisAddr string, redisChannel string) {
	rClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	incomingChannel := rClient.Subscribe(redisChannel)
	defer incomingChannel.Close()

	for {
		epmsg := event.Message{}

		msg, err := incomingChannel.ReceiveMessage(); 
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