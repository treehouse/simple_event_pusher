package client

import (
	"github.com/go-redis/redis"
	event "github.com/treehouse/simple_event_pusher/pkg/event"
)

// Stock redis events you can use in redis-cli:
// PUBLISH REDIS_PUBSUB_CHANNEL_PREFIX:channel-name '{ "event": "", "channel": "channel-name", "data": [{ "this": "is" }, "data"] }'
// PUBLISH event_pusher:hello-world '{ "event": "", "channel": "hello-world", "data": "this is data"}'
// PUBLISH event_pusher:hello-multiplex '{ "event": "", "channel": "hello-multiplex", "data": "this is data"}'
type RedisMessage redis.Message

func NewRedisMsg(channel, data string) event.Message {
	return &RedisMessage{
		Channel: channel,
		Payload: data,
	}
}

func (r RedisMessage) GetChannel() string { return r.Channel[len("event_pusher")+1:] }
func (r RedisMessage) Id() string         { return "" }
func (r RedisMessage) Event() string      { return "" }
func (r RedisMessage) Data() string       { return r.Payload }

type Subscription struct {
	pubsub *redis.PubSub
}

func (s Subscription) Close() error {
	return s.Close()
}

func (s Subscription) ReceiveMessage() (event.Message, error) {
	msg, err := s.pubsub.ReceiveMessage()
	msi := RedisMessage(*msg)
	return msi, err
}

func Redis(redisAddr string, redisChannel string) Incoming {
	rClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pubsub := rClient.PSubscribe(redisChannel + ":*")
	return Subscription{pubsub: pubsub}
}
