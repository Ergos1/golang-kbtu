package main

import (
	"context"
	"example.com/internal/store/psql/store/postgres"

	"example.com/internal/http"
)

func main() {
	store := postgres.NewDB()
	if err := store.Connect(); err != nil {
		panic(err)
	}
	defer store.Close()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}
