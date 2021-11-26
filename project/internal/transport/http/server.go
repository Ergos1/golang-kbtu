package http

import (
	"context"
	"log"
	"net/http"
	"time"

	messagebroker "example.com/internal/message_broker"
	"example.com/internal/store"
	"example.com/internal/transport/http/handlers"
	"example.com/pkg/cache/redis"

	"github.com/go-chi/chi"
)

type Server struct {
	Address string

	cache       *redis.RedisCache
	broker      messagebroker.MessageBroker
	store       store.Store
	idleConnsCh chan struct{}
	ctx         context.Context
}

type Response struct {
	Status  int
	Message string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Mount("/wallets", handlers.NewWalletHandler(s.store, s.broker, s.cache).Routes())

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] server runing on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] shutdown with error %v", err)
	}

	log.Println("[HTTP] shutdowned")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	// block before write or close
	<-s.idleConnsCh
}
