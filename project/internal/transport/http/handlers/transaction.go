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

type TransactionHandler struct {
	store  store.Store
	broker messagebroker.MessageBroker
	cache  *redis.RedisCache
}

func NewTransactionHandler(store store.Store, broker messagebroker.MessageBroker, cache *redis.RedisCache) *TransactionHandler {
	return &TransactionHandler{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (ah *TransactionHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ah.AllTransactions)
	r.Post("/", ah.CreateTransaction)
	r.Put("/", ah.UpdateTransaction)
	r.Delete("/{id}", ah.DeleteTransaction)
	r.Get("/{id}", ah.ByID)
	return r
}

func (ah *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := new(models.Transaction)
	if err := json.NewDecoder(r.Body).Decode(transaction); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Transactions().Create(r.Context(), transaction); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	ah.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории
	w.WriteHeader(http.StatusCreated)
}

func (ah *TransactionHandler) AllTransactions(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.TransactionFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		searchQuery = "transaction-" + searchQuery
		transactionFromCache := new(models.Transaction)
		if err := ah.cache.Get(r.Context(), searchQuery, transactionFromCache); err == nil {
			render.JSON(w, r, transactionFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	transactions, err := ah.store.Transactions().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		ah.cache.Set(r.Context(), searchQuery, transactions)
	}
	render.JSON(w, r, transactions)
}

func (ah *TransactionHandler) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	transactionFromCache := new(models.Transaction)
	if err := ah.cache.Get(r.Context(), id, transactionFromCache); err == nil {
		render.JSON(w, r, transactionFromCache)
		return
	}

	transaction, err := ah.store.Transactions().ByID(r.Context(), uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.cache.Set(r.Context(), id, transaction)
	render.JSON(w, r, transaction)
}

func (ah *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := new(models.Transaction)
	if err := json.NewDecoder(r.Body).Decode(transaction); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		transaction,
		validation.Field(&transaction.Id, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Transactions().Update(r.Context(), transaction); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.broker.Cache().Remove(transaction.Id)
}

func (ah *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Transactions().Delete(r.Context(), uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.broker.Cache().Remove(id)
}
