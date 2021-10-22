package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/internal/models"
	"example.com/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	Address string

	store       store.Store
	idleConnsCh chan struct{}
	ctx         context.Context
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		Address:     address,
		idleConnsCh: make(chan struct{}),
		store:       store,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Post("/nfts", func(w http.ResponseWriter, r *http.Request) {
		nft := new(models.NFT)
		nft.TokenURI = make(map[uint64]string)
		if err := json.NewDecoder(r.Body).Decode(nft); err != nil {
			fmt.Fprintf(w, "unknown error: %v", err)
			return
		}

		s.store.Create(r.Context(), nft)
	})

	r.Get("/nfts", func(w http.ResponseWriter, r *http.Request) {
		nfts, err := s.store.All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		render.JSON(w, r, nfts)
	})

	r.Get("/nfts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		nft, err := s.store.ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, nft)
	})

	r.Put("/nfts", func(w http.ResponseWriter, r *http.Request) {
		nft := new(models.NFT)
		if err := json.NewDecoder(r.Body).Decode(nft); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Update(r.Context(), nft)
	})

	r.Delete("/nfts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Delete(r.Context(), id)
	})

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
		return
	}

	log.Println("[HTTP] shutdowned")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	// block before write or close
	<-s.idleConnsCh
}
