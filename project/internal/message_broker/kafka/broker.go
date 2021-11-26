package kafka

import (
	"context"

	messagebroker "example.com/internal/message_broker"
	"example.com/pkg/cache/redis"
)

type Broker struct {
	brokers  []string
	clientID string

	cacheBroker messagebroker.CacheBroker
	cache       *redis.RedisCache
}

func NewBroker(brokers []string, cache *redis.RedisCache, clientID string) messagebroker.MessageBroker {
	return &Broker{
		brokers:  brokers,
		cache:    cache,
		clientID: clientID,
	}
}

func (b *Broker) Connect(ctx context.Context) error {
	brokers := []messagebroker.BrokerWithClient{b.Cache()}
	for _, broker := range brokers {
		if err := broker.Connect(ctx, b.brokers); err != nil {
			return err
		}
	}
	return nil
}

func (b *Broker) Close() error {
	brokers := []messagebroker.BrokerWithClient{b.Cache()}
	for _, broker := range brokers {
		if err := broker.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Broker) Cache() messagebroker.CacheBroker {
	if b.cacheBroker == nil {
		b.cacheBroker = NewCacheBroker(b.cache, b.clientID)
	}

	return b.cacheBroker
}
