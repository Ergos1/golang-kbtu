package http

import (
	"context"
	"encoding/json"
	"example.com/internal/store/psql/models"
	"example.com/internal/store/psql/store"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation/v4"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	Address string

	store store.Store
	idleConnsCh chan struct{}
	ctx         context.Context
}

type Response struct {
	Status  int
	Message string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		Address:     address,
		idleConnsCh: make(chan struct{}),
		store:       store,
	}
}

func collectionHandler(r *chi.Mux, s *Server) {
	r.Get("/collections", func(w http.ResponseWriter, r *http.Request) {
		collections, err := s.store.Collections().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, collections)
	})
	
	r.Post("/collections", func(w http.ResponseWriter, r *http.Request) {
		collection := new(models.Collection)
		if err := json.NewDecoder(r.Body).Decode(collection); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Collections().Create(r.Context(), collection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		collection, err := s.store.Collections().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, collection)
	})

	r.Put("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		collection := new(models.Collection)
		if err := json.NewDecoder(r.Body).Decode(collection); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		err := validation.ValidateStruct(
			collection,
			validation.Field(&collection.Id, validation.Required),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Collections().Update(r.Context(), collection); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	r.Delete("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Collections().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
}

func userHandler(r *chi.Mux, s *Server) {
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.Users().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, users)
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.Client)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Users().Create(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		user, err := s.store.Users().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, user)
	})

	r.Put("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.Client)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		err := validation.ValidateStruct(
			user,
			validation.Field(&user.Id, validation.Required),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Users().Update(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Users().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
}


func walletHandler(r *chi.Mux, s *Server) {
	r.Get("/wallets", func(w http.ResponseWriter, r *http.Request) {
		wallets, err := s.store.Wallets().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, wallets)
	})

	r.Post("/wallets", func(w http.ResponseWriter, r *http.Request) {
		wallet := new(models.Wallet)
		if err := json.NewDecoder(r.Body).Decode(wallet); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Wallets().Create(r.Context(), wallet); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/wallets/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		wallet, err := s.store.Wallets().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, wallet)
	})

	r.Put("/wallets/{id}", func(w http.ResponseWriter, r *http.Request) {
		wallet := new(models.Wallet)
		if err := json.NewDecoder(r.Body).Decode(wallet); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		err := validation.ValidateStruct(
			wallet,
			validation.Field(&wallet.Id, validation.Required),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Wallets().Update(r.Context(), wallet); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	r.Delete("/wallets/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Wallets().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
}

func transactionHandler(r *chi.Mux, s *Server) {
	r.Get("/transactions", func(w http.ResponseWriter, r *http.Request) {
		rows, err := s.db.Queryx("SELECT * FROM TRANSACTIOn")
		if err != nil {
			render.JSON(w, r, Response{
				Status:  500,
				Message: "SERVER's DB ERROR",
			})
			return
		}
		var transactions []models.Transaction
		transaction := new(models.Transaction)
		for rows.Next(){
			err = rows.StructScan(transaction)
			transactions = append(transactions, *transaction)
		}
		render.JSON(w, r, transactions)
	})

	r.Post("/transactions", func(w http.ResponseWriter, r *http.Request) {
		transaction := new(models.Transaction)
		if err := json.NewDecoder(r.Body).Decode(transaction); err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "Fields of collection are incorrect",
			})
			return
		}
		_, err := s.db.NamedExec(`INSERT INTO TRANSACTION(walletid, touserid, amount, description)
										values(:walletid, :touserid, :amount, :description)`, transaction)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "ASDASDAS",
			})
			return
		}
		render.JSON(w, r, Response{
			Status:  201,
			Message: "Created collection",
		})
	})

	r.Get("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}
		transaction := new(models.Transaction)
		err = s.db.Get(transaction, "SELECT * FROM TRANSACTION WHERE id=$1", id)
		if err != nil || *transaction == (models.Transaction{}){
			render.JSON(w, r, Response{
				Status:  404,
				Message: "NOT FOUND",
			})
			return
		}
		render.JSON(w, r, *transaction)
	})

	r.Put("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		newTransaction := new(models.Transaction)
		if err = json.NewDecoder(r.Body).Decode(newTransaction); err != nil {
			render.JSON(w, r, Response{
				Status: 400,
				Message: "BAD FIELDS",
			})
			return
		}
		var query []string
		v := reflect.ValueOf(*newTransaction)
		typeOf := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Interface() == reflect.Zero(reflect.TypeOf(v.Field(i).Interface())).Interface(){
				continue
			}
			query = append(query, fmt.Sprintf("%s=%v", typeOf.Field(i).Name, v.Field(i).Interface()))
		}
		//fmt.Println(query, id)
		//res := fmt.Sprintf(
		//	`UPDATE Collection SET %s where id=%d`, strings.Join(query, ", "), id)
		//fmt.Println(res)
		_, err = s.db.Exec(fmt.Sprintf(
			`UPDATE Transaction SET %s where id=%d`, strings.Join(query, ", "), id))
		render.JSON(w, r, Response{
			Status:  200,
			Message: "Success changed",
		})
	})
	r.Delete("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		if res, err := s.db.Exec(fmt.Sprintf("DELETE FROM transaction where id=%d", id)); err != nil  {
			render.JSON(w, r, Response{
				Status:  404,
				Message: "NOT FOUND TO DELETE",
			})
			return
		} else if count, err := res.RowsAffected(); count == 0 || err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "Not exist",
			})
			return
		}
		render.JSON(w, r, Response{
			Status:  200,
			Message: "DELETED",
		})
	})
}

func nftHandler(r *chi.Mux, s *Server) {
	r.Get("/nfts", func(w http.ResponseWriter, r *http.Request) {
		rows, err := s.db.Queryx("SELECT * FROM NonFungibleToken")
		if err != nil {
			render.JSON(w, r, Response{
				Status:  500,
				Message: "SERVER's DB ERROR",
			})
			return
		}
		var nfts []models.NonFungibleToken
		nft := new(models.NonFungibleToken)
		for rows.Next(){
			err = rows.StructScan(nft)
			nfts = append(nfts, *nft)
		}
		render.JSON(w, r, nfts)
	})

	r.Post("/nfts", func(w http.ResponseWriter, r *http.Request) {
		nft := new(models.NonFungibleToken)
		if err := json.NewDecoder(r.Body).Decode(nft); err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "Fields of collection are incorrect",
			})
			return
		}
		_, err := s.db.NamedExec(`INSERT INTO NonFungibleToken(likes, collectionid, ownerid,
                            price, royalties, title, description) values(:likes, :collectionid, :ownerid,
                            :price, :royalties, :title, :description)`, nft)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "ASDASDAS",
			})
			return
		}
		render.JSON(w, r, Response{
			Status:  201,
			Message: "Created collection",
		})
	})

	r.Get("/nfts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}
		nft := new(models.NonFungibleToken)
		err = s.db.Get(nft, "SELECT * FROM NonFungibleToken WHERE id=$1", id)
		if err != nil || *nft == (models.NonFungibleToken{}){
			render.JSON(w, r, Response{
				Status:  404,
				Message: "NOT FOUND",
			})
			return
		}
		render.JSON(w, r, *nft)
	})

	r.Put("/nfts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		newNft := new(models.NonFungibleToken)
		if err = json.NewDecoder(r.Body).Decode(newNft); err != nil {
			render.JSON(w, r, Response{
				Status: 400,
				Message: "BAD FIELDS",
			})
			return
		}
		var query []string
		v := reflect.ValueOf(*newNft)
		typeOf := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Interface() == reflect.Zero(reflect.TypeOf(v.Field(i).Interface())).Interface(){
				continue
			}
			query = append(query, fmt.Sprintf("%s=%v", typeOf.Field(i).Name, v.Field(i).Interface()))
		}
		//fmt.Println(query, id)
		//res := fmt.Sprintf(
		//	`UPDATE Collection SET %s where id=%d`, strings.Join(query, ", "), id)
		fmt.Println(query, id)
		_, err = s.db.Exec(fmt.Sprintf(
			`UPDATE NonFungibleToken SET %s where id=%d`, strings.Join(query, ", "), id))
		render.JSON(w, r, Response{
			Status:  200,
			Message: "Success changed",
		})
	})
	r.Delete("/nfts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		if res, err := s.db.Exec(fmt.Sprintf("DELETE FROM NonFungibleToken where id=%d", id)); err != nil  {
			render.JSON(w, r, Response{
				Status:  404,
				Message: "NOT FOUND TO DELETE",
			})
			return
		} else if count, err := res.RowsAffected(); count == 0 || err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "Not exist",
			})
			return
		}
		render.JSON(w, r, Response{
			Status:  200,
			Message: "DELETED",
		})
	})
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	collectionHandler(r, s)
	walletHandler(r, s)
	userHandler(r, s)
	transactionHandler(r, s)
	nftHandler(r, s)
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
