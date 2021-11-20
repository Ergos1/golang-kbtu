package main

import (
	"context"
	"example.com/internal/config"
	"example.com/internal/store/psql"
	http2 "example.com/internal/transport/http"
	cache2 "example.com/pkg/cache/redis"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	_ = config.NewConfig()
	store := psql.NewDB()
	if err := store.Connect(); err != nil {
		panic(err)
	}
	defer store.Close()
	cache := cache2.NewRedisCache()
	cache.Connect("localhost:6379", 0, 0)
	srv := http2.NewServer(context.Background(),
		http2.WithCache(cache),
		http2.WithStore(store),
		http2.WithAddress(":8080"))
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}
