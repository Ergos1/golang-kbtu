package main

import (
	"context"
	"example.com/internal/store/psql"
	"log"

	"example.com/internal/http"
)

func main() {
	db := psql.NewDB()
	defer db.Close()

	srv := http.NewServer(context.Background(), ":8080", db)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
