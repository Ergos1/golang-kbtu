package http

import (
	"context"
	"encoding/json"
	"example.com/internal/store/psql/models"
	"fmt"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
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

	db       *sqlx.DB
	idleConnsCh chan struct{}
	ctx         context.Context
}

type Response struct {
	Status  int
	Message string
}

func NewServer(ctx context.Context, address string, db *sqlx.DB) *Server {
	return &Server{
		ctx:         ctx,
		Address:     address,
		idleConnsCh: make(chan struct{}),
		db:       db,
	}
}

func collectionHandler(r *chi.Mux, s *Server) {
	r.Get("/collections", func(w http.ResponseWriter, r *http.Request) {
		rows, err := s.db.Queryx("SELECT * FROM COLLECTION")
		if err != nil {
			render.JSON(w, r, Response{
				Status:  500,
				Message: "SERVER's DB ERROR",
			})
			return
		}
		var collections []models.Collection
		collection := new(models.Collection)
		for rows.Next(){
			err = rows.StructScan(collection)
			collections = append(collections, *collection)
		}
		render.JSON(w, r, collections)
	})
	
	r.Post("/collections", func(w http.ResponseWriter, r *http.Request) {
		collection := new(models.Collection)
		collection.SetDefaultId()
		if err := json.NewDecoder(r.Body).Decode(collection); err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "Fields of collection are incorrect",
			})
			return
		}
		_, err := s.db.NamedExec(`INSERT INTO Collection(name, symbol, description, ownerid)
								values(:name, :symbol, :description, :ownerid)`, collection)
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

	r.Get("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}
		collection := new(models.Collection)
		err = s.db.Get(collection, "SELECT * FROM Collection WHERE id=$1", id)
		if err != nil || *collection == (models.Collection{}){
			render.JSON(w, r, Response{
				Status:  404,
				Message: "NOT FOUND",
			})
			return
		}
		render.JSON(w, r, *collection)
	})

	r.Put("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		newCollection := new(models.Collection)
		if err = json.NewDecoder(r.Body).Decode(newCollection); err != nil {
			render.JSON(w, r, Response{
				Status: 400,
				Message: "BAD FIELDS",
			})
			return
		}
		var query []string
		v := reflect.ValueOf(*newCollection)
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
			`UPDATE Collection SET %s where id=%d`, strings.Join(query, ", "), id))
		render.JSON(w, r, Response{
			Status:  200,
			Message: "Success changed",
		})
	})
	r.Delete("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		if res, err := s.db.Exec(fmt.Sprintf("DELETE FROM collection where id=%d", id)); err != nil  {
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

func userHandler(r *chi.Mux, s *Server) {
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		rows, err := s.db.Queryx("SELECT * FROM Client")
		if err != nil {
			render.JSON(w, r, Response{
				Status:  500,
				Message: "SERVER's DB ERROR",
			})
			return
		}
		var users []models.Client
		user := new(models.Client)
		for rows.Next(){
			err = rows.StructScan(user)
			users = append(users, *user)
			fmt.Println(*user)
		}
		render.JSON(w, r, users)
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.Client)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "Fields of collection are incorrect",
			})
			return
		}
		_, err := s.db.NamedExec(`INSERT INTO Client(walletid,username,email,password) 
										values(:walletid,:username, :email,:password)`, user)
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

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}
		user := new(models.Client)
		err = s.db.Get(user, "SELECT * FROM Client WHERE id=$1", id)
		if err != nil || *user == (models.Client{}){
			render.JSON(w, r, Response{
				Status:  404,
				Message: "NOT FOUND",
			})
			return
		}
		render.JSON(w, r, *user)
	})

	r.Put("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		newUser := new(models.Client)
		if err = json.NewDecoder(r.Body).Decode(newUser); err != nil {
			render.JSON(w, r, Response{
				Status: 400,
				Message: "BAD FIELDS",
			})
			return
		}
		var query []string
		v := reflect.ValueOf(*newUser)
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
			`UPDATE Client SET %s where id=%d`, strings.Join(query, ", "), id))
		render.JSON(w, r, Response{
			Status:  200,
			Message: "Success changed",
		})
	})
	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		if res, err := s.db.Exec(fmt.Sprintf("DELETE FROM client where id=%d", id)); err != nil  {
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


func walletHandler(r *chi.Mux, s *Server) {
	r.Get("/wallets", func(w http.ResponseWriter, r *http.Request) {
		rows, err := s.db.Queryx("SELECT * FROM Wallet")
		if err != nil {
			render.JSON(w, r, Response{
				Status:  500,
				Message: "SERVER's DB ERROR",
			})
			return
		}
		var wallets []models.Wallet
		wallet := new(models.Wallet)
		for rows.Next(){
			err = rows.StructScan(wallet)
			wallets = append(wallets, *wallet)
		}
		render.JSON(w, r, wallets)
	})

	r.Post("/wallets", func(w http.ResponseWriter, r *http.Request) {
		wallet := new(models.Wallet)
		if err := json.NewDecoder(r.Body).Decode(wallet); err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "Fields of collection are incorrect",
			})
			return
		}
		_, err := s.db.NamedExec(`INSERT INTO Wallet(balance) values(:balance)`, wallet)
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

	r.Get("/wallets/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}
		wallet := new(models.Wallet)
		err = s.db.Get(wallet, "SELECT * FROM Wallet WHERE id=$1", id)
		if err != nil || *wallet == (models.Wallet{}){
			render.JSON(w, r, Response{
				Status:  404,
				Message: "NOT FOUND",
			})
			return
		}
		render.JSON(w, r, *wallet)
	})

	r.Put("/wallets/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		newWallet := new(models.Wallet)
		if err = json.NewDecoder(r.Body).Decode(newWallet); err != nil {
			render.JSON(w, r, Response{
				Status: 400,
				Message: "BAD FIELDS",
			})
			return
		}
		var query []string
		v := reflect.ValueOf(*newWallet)
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
			`UPDATE Wallet SET %s where id=%d`, strings.Join(query, ", "), id))
		render.JSON(w, r, Response{
			Status:  200,
			Message: "Success changed",
		})
	})
	r.Delete("/wallets/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, Response{
				Status:  400,
				Message: "BAD ID",
			})
			return
		}

		if res, err := s.db.Exec(fmt.Sprintf("DELETE FROM wallet where id=%d", id)); err != nil  {
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