package main

import (
	"example.com/api"
	grpc2 "example.com/internal/store/inmemory/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":5000"
)

func main() {
	store := grpc2.NewDB()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot listen to %s: %v", port, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	collectionService := store.Collections()
	nonFungibleTokenService := store.NonFungibleToken()
	userService := store.User()
	walletService := store.Wallet()

	api.RegisterCollectionServiceServer(grpcServer, collectionService)
	api.RegisterNonFungibleTokenServiceServer(grpcServer, nonFungibleTokenService)
	api.RegisterUserServiceServer(grpcServer, userService)
	api.RegisterWalletServiceServer(grpcServer, walletService)

	log.Printf("Serving on %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve on %v: %v", listener.Addr(), err)
	}
}