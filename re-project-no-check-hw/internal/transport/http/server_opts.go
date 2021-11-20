package http

import (
	"example.com/internal/database/psql/store"
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

func WithCache(cache redis.Cache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}