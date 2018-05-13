package client

import (
	"github.com/go-redis/redis"
)

type RedisMessage redis.Message

func (r RedisMessage) GetChannel() string {
	return r.Channel
}

func (r RedisMessage) GetPayload() string {
	return r.Payload
}

type Subscription struct {
	pubsub *redis.PubSub
}

func (s Subscription) Close() error {
	return s.Close()
}

func (s Subscription) ReceiveMessage() (MessageInterface, error) {
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
	pubsub := rClient.Subscribe(redisChannel)
	return Subscription{ pubsub: pubsub }
}

