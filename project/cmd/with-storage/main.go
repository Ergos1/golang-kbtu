package main

import (
	"context"
	"log"

	"example.com/internal/http"
	"example.com/internal/store/inmemory"
)

func main() {
	store := inmemory.NewDB()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
