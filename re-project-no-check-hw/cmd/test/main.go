package main

import (
	"context"
	"fmt"

	"example.com/pkg/cache/redis"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	cache := redis.NewRedisCache()
	cache.Connect("localhost:6379", 0, 0)
	user := new(User)
	// user.Name = "VASIA"
	// user.Id = 1
	// p, _ := json.Marshal(user)
	// cache.Set(context.Background(), "123", p)
	cache.Get(context.Background(), "123", user)
	fmt.Print(user)
}
