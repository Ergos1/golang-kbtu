package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	messagebroker "example.com/internal/message_broker"
	"example.com/internal/models"
	"example.com/internal/store"
	"example.com/pkg/cache/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
)

type WalletHandler struct {
	store  store.Store
	broker messagebroker.MessageBroker
	cache  *redis.RedisCache
}

func NewWalletHandler(store store.Store, broker messagebroker.MessageBroker, cache *redis.RedisCache) *WalletHandler {
	return &WalletHandler{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (wh *WalletHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", wh.AllWallets)
	r.Post("/", wh.CreateWalet)
	r.Put("/", wh.UpdateWallet)
	r.Delete("/{id}", wh.DeleteWallet)
	r.Get("/{id}", wh.ByID)
	return r
	// r.Post("/", )
}

func (wh *WalletHandler) CreateWalet(w http.ResponseWriter, r *http.Request) {
	wallet := new(models.Wallet)
	if err := json.NewDecoder(r.Body).Decode(wallet); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := wh.store.Wallets().Create(r.Context(), wallet); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	wh.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории
	w.WriteHeader(http.StatusCreated)
}

func (wh *WalletHandler) AllWallets(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.WalletFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		searchQuery = "wallet-" + searchQuery
		walletFromCache := new([]*models.Wallet)
		if err := wh.cache.Get(r.Context(), searchQuery, walletFromCache); err == nil {
			// fmt.Println("FROM CACHE")
			render.JSON(w, r, walletFromCache)
			return
		}
		fmt.Println(walletFromCache)

		filter.Query = &searchQuery
	}

	wallets, err := wh.store.Wallets().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		wh.cache.Set(r.Context(), searchQuery, wallets)
	}
	render.JSON(w, r, wallets)
}

func (wh *WalletHandler) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	walletFromCache := new(models.Wallet)
	if err := wh.cache.Get(r.Context(), id, walletFromCache); err == nil {
		render.JSON(w, r, walletFromCache)
		return
	}

	wallet, err := wh.store.Wallets().ByID(r.Context(), uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	wh.cache.Set(r.Context(), id, wallet)
	render.JSON(w, r, wallet)
}

func (wh *WalletHandler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
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

	if err := wh.store.Wallets().Update(r.Context(), wallet); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	wh.broker.Cache().Remove(wallet.Id)
}

func (wh *WalletHandler) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := wh.store.Wallets().Delete(r.Context(), uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	wh.broker.Cache().Remove(id)
}
