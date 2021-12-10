package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func (cache *RedisCache) Connect(host string, db int, expires time.Duration) {
	cache.client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       db,
	})
}

func (cache *RedisCache) Close() error {
	err := cache.client.Close()
	return err
}

func (cache *RedisCache) Get(ctx context.Context, key interface{}, dest interface{}) error {
	keyStr := fmt.Sprintf("%v", key)
	result, err := cache.client.Get(ctx, keyStr).Result()
	fmt.Println(result)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (cache *RedisCache) Set(ctx context.Context, key interface{}, value interface{}) error {
	keyStr := fmt.Sprintf("%v", key)
	valueParsed, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = cache.client.Set(ctx, keyStr, valueParsed, 0).Err()
	fmt.Println(err)
	return err
}

// Idea of Remove is like:
// I get like wallet-... || account-...
// After I just create struct by new(...)
// Save in temp field after remove from cache
// Remove by id item and again set it
// That's all
func (cache *RedisCache) Remove(ctx context.Context, key interface{}) error {
	// info := strings.Split(fmt.Sprintf("%v", key), "-")
	// switch info[0] {
	// case "wallet":
	// 	fmt.Print("Implement me EPA DONT FORGET")
	// case "account":
	// 	fmt.Print("Implement me EPA DONT FORGET")
	// case "...":
	// 	fmt.Print("Implement me EPA DONT FORGET")
	// default:
	// 	fmt.Print("Implement me EPA DONT FORGET")
	// }

	// var oldResult interface{}
	// cache.Get(ctx, "all", oldResult)

	// fmt.Println(oldResult) // Pls deadline gonna kill me, give me more time, I promise I will do it

	result := cache.client.Del(ctx, fmt.Sprintf("%v", key))
	if result.Err() != nil {
		return result.Err()
	}

	return nil
	// return cache.Purge(ctx)
}

func (cache *RedisCache) Purge(ctx context.Context) error {
	result := cache.client.FlushAll(ctx)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func NewRedisCache() *RedisCache {
	return &RedisCache{}
}
