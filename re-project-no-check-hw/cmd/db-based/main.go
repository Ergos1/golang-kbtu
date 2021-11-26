package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"example.com/internal/config"
	"example.com/internal/store/psql"
	"example.com/internal/transport/http"
	redis "example.com/pkg/cache/redis"
	"github.com/joho/godotenv"
	"example.com/internal/message_broker/kafka"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	cfg := config.NewConfig()

	store := psql.NewDB()
	fmt.Println(cfg.Database.Uri())
	if err := store.Connect(cfg.Database.Uri()); err != nil {
		panic(err)
	}
	defer store.Close()
	
	cache := redis.NewRedisCache()
	cache.Connect(cfg.Redis.Host, cfg.Redis.Db, cfg.Redis.Expires)
	cache.Purge(ctx) // Delete if not need purge cache 

	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer3")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	srv := http.NewServer(
		ctx,
		http.WithCache(cache),
		http.WithStore(store),
		http.WithAddress(":8080"),
		http.WithBroker(broker))
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[warning] Caught termination signal")
	cancel()
}