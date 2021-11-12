package store_grpc

import (
	"example.com/api"
)

type Store interface {
	Collections() CollectionRepository
	NonFungibleToken() NonFungibleTokenRepository
	User() UserRepository
	Wallet() WalletRepository
}

type CollectionRepository interface {
	api.CollectionServiceServer
}

type NonFungibleTokenRepository interface {
	api.NonFungibleTokenServiceServer
}

type UserRepository interface {
	api.UserServiceServer
}

type WalletRepository interface {
	api.WalletServiceServer
}

