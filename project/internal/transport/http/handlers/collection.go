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

type CollectionHandler struct {
	store  store.Store
	broker messagebroker.MessageBroker
	cache  *redis.RedisCache
}

func NewCollectionHandler(store store.Store, broker messagebroker.MessageBroker, cache *redis.RedisCache) *CollectionHandler {
	return &CollectionHandler{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (ah *CollectionHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ah.AllCollections)
	r.Post("/", ah.CreateCollection)
	r.Put("/", ah.UpdateCollection)
	r.Delete("/{id}", ah.DeleteCollection)
	r.Get("/{id}", ah.ByID)
	return r
	// r.Post("/", )
}

func (ah *CollectionHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	collection := new(models.Collection)
	if err := json.NewDecoder(r.Body).Decode(collection); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Collections().Create(r.Context(), collection); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	ah.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории
	w.WriteHeader(http.StatusCreated)
}

func (ah *CollectionHandler) AllCollections(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.CollectionFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		searchQuery = "collection-" + searchQuery
		collectionFromCache := new(models.Collection)
		if err := ah.cache.Get(r.Context(), searchQuery, collectionFromCache); err == nil {
			render.JSON(w, r, collectionFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	collections, err := ah.store.Collections().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		ah.cache.Set(r.Context(), searchQuery, collections)
	}
	render.JSON(w, r, collections)
}

func (ah *CollectionHandler) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	collectionFromCache := new(models.Collection)
	if err := ah.cache.Get(r.Context(), id, collectionFromCache); err == nil {
		render.JSON(w, r, collectionFromCache)
		return
	}

	collection, err := ah.store.Collections().ByID(r.Context(), uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.cache.Set(r.Context(), id, collection)
	render.JSON(w, r, collection)
}

func (ah *CollectionHandler) UpdateCollection(w http.ResponseWriter, r *http.Request) {
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

	if err := ah.store.Collections().Update(r.Context(), collection); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.broker.Cache().Remove(collection.Id)
}

func (ah *CollectionHandler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Collections().Delete(r.Context(), uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.broker.Cache().Remove(id)
}
