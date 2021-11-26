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

type AccountHandler struct {
	store  store.Store
	broker messagebroker.MessageBroker
	cache  *redis.RedisCache
}

func NewAccountHandler(store store.Store, broker messagebroker.MessageBroker, cache *redis.RedisCache) *AccountHandler {
	return &AccountHandler{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (ah *AccountHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ah.AllAccounts)
	r.Post("/", ah.CreateAccount)
	r.Put("/", ah.UpdateAccount)
	r.Delete("/{id}", ah.DeleteAccount)
	r.Get("/{id}", ah.ByID)
	return r
	// r.Post("/", )
}

func (ah *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := new(models.Account)
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Accounts().Create(r.Context(), account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	ah.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории
	w.WriteHeader(http.StatusCreated)
}

func (ah *AccountHandler) AllAccounts(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.AccountFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		searchQuery = "account-" + searchQuery
		accountFromCache := new(models.Account)
		if err := ah.cache.Get(r.Context(), searchQuery, accountFromCache); err == nil {
			render.JSON(w, r, accountFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	accounts, err := ah.store.Accounts().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		ah.cache.Set(r.Context(), searchQuery, accounts)
	}
	render.JSON(w, r, accounts)
}

func (ah *AccountHandler) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	accountFromCache := new(models.Account)
	if err := ah.cache.Get(r.Context(), id, accountFromCache); err == nil {
		render.JSON(w, r, accountFromCache)
		return
	}

	account, err := ah.store.Accounts().ByID(r.Context(), uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.cache.Set(r.Context(), id, account)
	render.JSON(w, r, account)
}

func (ah *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	account := new(models.Account)
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		account,
		validation.Field(&account.Id, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Accounts().Update(r.Context(), account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.broker.Cache().Remove(account.Id)
}

func (ah *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Accounts().Delete(r.Context(), uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// ah.broker.Cache().Remove(id)
	ah.broker.Cache().Purge()
}
