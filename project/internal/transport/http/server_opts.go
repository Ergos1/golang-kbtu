package http

import (
	messagebroker "example.com/internal/message_broker"
	"example.com/internal/store"
	"example.com/pkg/cache/redis"
)

type ServerOption func(srv *Server)

func WithAddress(address string) ServerOption {
	return func(srv *Server) {
		srv.Address = address
	}
}

func WithStore(store store.Store) ServerOption {
	return func(srv *Server) {
		srv.store = store
	}
}

func WithCache(cache *redis.RedisCache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}

func WithBroker(broker messagebroker.MessageBroker) ServerOption {
	return func(srv *Server) {
		srv.broker = broker
	}
}
