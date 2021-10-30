package grpc

import (
	"sync"

	"example.com/api"

	"example.com/internal/store"
)

type DB struct {
	collectionRepository       api.CollectionServiceServer
	nonFungibleTokenRepository api.NonFungibleTokenServiceServer
	walletRepository           api.WalletServiceServer
	userRepository             api.UserServiceServer

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Collections() store.CollectionRepository {
	if db.collectionRepository == nil {
		db.collectionRepository = &CollectionRepo{
			data: make(map[uint64]*api.Collection),
			mu:   new(sync.RWMutex),
		}
	}
	return db.collectionRepository
}

func (db *DB) NonFungibleToken() store.NonFungibleTokenRepository {
	if db.nonFungibleTokenRepository == nil {
		db.nonFungibleTokenRepository = &NonFungibleTokenRepo{
			data: make(map[uint64]*api.NonFungibleToken),
			mu:   new(sync.RWMutex),
		}
	}
	return db.nonFungibleTokenRepository
}

func (db *DB) User() store.UserRepository {
	if db.userRepository == nil {
		db.userRepository = &UserRepo{
			data: make(map[uint64]*api.User),
			mu:   new(sync.RWMutex),
		}
	}
	return db.userRepository
}

func (db *DB) Wallet() store.WalletRepository {
	if db.walletRepository == nil {
		db.walletRepository = &WalletRepo{
			data: make(map[uint64]*api.Wallet),
			mu:   new(sync.RWMutex),
		}
	}
	return db.walletRepository
}
