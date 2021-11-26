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

type CategoryHandler struct {
	store  store.Store
	broker messagebroker.MessageBroker
	cache  *redis.RedisCache
}

func NewCategoryHandler(store store.Store, broker messagebroker.MessageBroker, cache *redis.RedisCache) *CategoryHandler {
	return &CategoryHandler{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (ah *CategoryHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ah.AllCategories)
	r.Post("/", ah.CreateCategory)
	r.Put("/", ah.UpdateCategory)
	r.Delete("/{id}", ah.DeleteCategory)
	r.Get("/{id}", ah.ByID)
	return r
	// r.Post("/", )
}

func (ah *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := new(models.Category)
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Categories().Create(r.Context(), category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	ah.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории
	w.WriteHeader(http.StatusCreated)
}

func (ah *CategoryHandler) AllCategories(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.CategoryFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		searchQuery = "category-" + searchQuery
		categoryFromCache := new(models.Category)
		if err := ah.cache.Get(r.Context(), searchQuery, categoryFromCache); err == nil {
			render.JSON(w, r, categoryFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	categories, err := ah.store.Categories().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		ah.cache.Set(r.Context(), searchQuery, categories)
	}
	render.JSON(w, r, categories)
}

func (ah *CategoryHandler) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	categoryFromCache := new(models.Category)
	if err := ah.cache.Get(r.Context(), id, categoryFromCache); err == nil {
		render.JSON(w, r, categoryFromCache)
		return
	}

	category, err := ah.store.Categories().ByID(r.Context(), uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.cache.Set(r.Context(), id, category)
	render.JSON(w, r, category)
}

func (ah *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	category := new(models.Category)
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		category,
		validation.Field(&category.Id, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Categories().Update(r.Context(), category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.broker.Cache().Remove(category.Id)
}

func (ah *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Categories().Delete(r.Context(), uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.broker.Cache().Remove(id)
}
