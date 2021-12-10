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

type AssetHandler struct {
	store  store.Store
	broker messagebroker.MessageBroker
	cache  *redis.RedisCache
}

func NewAssetHandler(store store.Store, broker messagebroker.MessageBroker, cache *redis.RedisCache) *AssetHandler {
	return &AssetHandler{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (ah *AssetHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ah.AllAssets)
	r.Post("/", ah.CreateAsset)
	r.Put("/", ah.UpdateAsset)
	r.Get("/{id}", ah.ByID)
	return r
	// r.Post("/", )
}

func (ah *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	asset := new(models.Asset)
	if err := json.NewDecoder(r.Body).Decode(asset); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Assets().Create(r.Context(), asset); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	ah.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории
	w.WriteHeader(http.StatusCreated)
}

func (ah *AssetHandler) AllAssets(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.AssetFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		searchQuery = "asset-" + searchQuery
		assetFromCache := new(models.Asset)
		if err := ah.cache.Get(r.Context(), searchQuery, assetFromCache); err == nil {
			render.JSON(w, r, assetFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	assets, err := ah.store.Assets().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		ah.cache.Set(r.Context(), searchQuery, assets)
	}
	render.JSON(w, r, assets)
}

func (ah *AssetHandler) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	assetFromCache := new(models.Asset)
	if err := ah.cache.Get(r.Context(), id, assetFromCache); err == nil {
		render.JSON(w, r, assetFromCache)
		return
	}

	asset, err := ah.store.Assets().ByID(r.Context(), uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.cache.Set(r.Context(), id, asset)
	render.JSON(w, r, asset)
}

func (ah *AssetHandler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	asset := new(models.Asset)
	if err := json.NewDecoder(r.Body).Decode(asset); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(
		asset,
		validation.Field(&asset.Id, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ah.store.Assets().Update(r.Context(), asset); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	ah.broker.Cache().Remove(asset.Id)
}

