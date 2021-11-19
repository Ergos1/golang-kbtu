package main

import (
	"context"
	"example.com/internal/store/psql/store/postgres"
	cache2 "example.com/pkg/cache"

	"example.com/internal/http"
)

func main() {
	store := postgres.NewDB()
	if err := store.Connect(); err != nil {
		panic(err)
	}
	defer store.Close()
	cache := cache2.NewRedisCache()
	cache.Connect("localhost:6379", 0, 0)
	srv := http.NewServer(context.Background(),
		http.WithCache(cache),
		http.WithStore(store),
		http.WithAddress(":8080"))
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}
