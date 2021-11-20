package redis

import (
	"context"
	json2 "encoding/json"
	"example.com/internal/database/psql/models"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Cache interface {
	Connect(host string, db int, expires time.Duration)

	Wallets() WalletsCache
	Users() UsersCache
}

type redisCache struct {
	client  *redis.Client
	users   UsersCache
	wallets WalletsCache
}

func (cache *redisCache) Connect(host string, db int, expires time.Duration) {
	cache.client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       db,
	})
}

func NewRedisCache() Cache {
	return &redisCache{}
}

func (cache *redisCache) Wallets() WalletsCache {
	if cache.wallets == nil {
		cache.wallets = NewWalletsCache(cache.client)
	}
	return cache.wallets
}

type WalletCache struct {
	client *redis.Client
}

func NewWalletsCache(client *redis.Client) WalletsCache {
	return &WalletCache{
		client: client,
	}
}

func (wc WalletCache) Set(key string, value *models.Wallets) {
	json, err := json2.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}

	wc.client.Set(context.Background(), key, json, 0)
}
func (wc WalletCache) Get(key string) *models.Wallets {
	value, err := wc.client.Get(context.Background(), key).Result()

	if err != nil {
		return nil
	}
	wallet := &models.Wallets{}
	err = json2.Unmarshal([]byte(value), wallet)
	if err != nil {
		return nil
	}

	return wallet
}

func (wc WalletCache) Delete(key string) {
	if _, err := wc.client.Get(context.Background(), key).Result(); err == nil {
		wc.client.Del(context.Background(), key)
	}
}

func (cache *redisCache) Users() UsersCache {
	if cache.users == nil {
		cache.users = NewUsersCache(cache.client)
	}
	return cache.users
}

type UserCache struct {
	client *redis.Client
}

func NewUsersCache(client *redis.Client) UsersCache {
	return &UserCache{
		client: client,
	}
}

func (u UserCache) Set(key string, value *models.Clients) {
	json, err := json2.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}

	u.client.Set(context.Background(), key, json, 0)
}
func (u UserCache) Get(key string) *models.Clients {
	value, err := u.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	client := &models.Clients{}
	err = json2.Unmarshal([]byte(value), client)
	if err != nil {
		return nil
	}

	return client
}


func (u UserCache) SetAll(key string, values []*models.Clients) {
	json, err := json2.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}

	u.client.Set(context.Background(), key, json, 0)
}

func (u UserCache) GetAll(key string) []*models.Clients {
	value, err := u.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	clients := make([]*models.Clients, 0)
	err = json2.Unmarshal([]byte(value), clients)
	if err != nil {
		return nil
	}

	return clients
}

func (u UserCache) Delete(key string) {
	if _, err := u.client.Get(context.Background(), key).Result(); err == nil {
		u.client.Del(context.Background(), key)
	}
}
