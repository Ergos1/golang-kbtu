package grpc

import (
	"example.com/internal/store/store-grpc"
	"sync"

	"example.com/api"
)

type DB struct {
	collectionRepository       api.CollectionServiceServer
	nonFungibleTokenRepository api.NonFungibleTokenServiceServer
	walletRepository           api.WalletServiceServer
	userRepository             api.UserServiceServer

	mu *sync.RWMutex
}

func NewDB() store_grpc.Store {
	return &DB{
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Collections() store_grpc.CollectionRepository {
	if db.collectionRepository == nil {
		db.collectionRepository = &CollectionRepo{
			data: make(map[uint64]*api.Collection),
			mu:   new(sync.RWMutex),
		}
	}
	return db.collectionRepository
}

func (db *DB) NonFungibleToken() store_grpc.NonFungibleTokenRepository {
	if db.nonFungibleTokenRepository == nil {
		db.nonFungibleTokenRepository = &NonFungibleTokenRepo{
			data: make(map[uint64]*api.NonFungibleToken),
			mu:   new(sync.RWMutex),
		}
	}
	return db.nonFungibleTokenRepository
}

func (db *DB) User() store_grpc.UserRepository {
	if db.userRepository == nil {
		db.userRepository = &UserRepo{
			data: make(map[uint64]*api.User),
			mu:   new(sync.RWMutex),
		}
	}
	return db.userRepository
}

func (db *DB) Wallet() store_grpc.WalletRepository {
	if db.walletRepository == nil {
		db.walletRepository = &WalletRepo{
			data: make(map[uint64]*api.Wallet),
			mu:   new(sync.RWMutex),
		}
	}
	return db.walletRepository
}
