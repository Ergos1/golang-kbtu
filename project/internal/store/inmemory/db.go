package inmemory

import (
	"context"
	"fmt"
	"sync"

	"example.com/internal/models"
	"example.com/internal/store"
)

type DB struct {
	data map[uint64]*models.NFT
	mu   *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[uint64]*models.NFT),
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, nft *models.NFT) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[nft.ID] = nft
	return nil
}

func (db *DB) All(ctx context.Context) ([]*models.NFT, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	nfts := make([]*models.NFT, 0, len(db.data))
	for _, nft := range db.data {
		nfts = append(nfts, nft)
	}

	return nfts, nil
}

func (db *DB) ByID(ctx context.Context, id uint64) (*models.NFT, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	nft, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("[error] no nft with id %d", id)
	}

	return nft, nil
}

func (db *DB) Update(ctx context.Context, nft *models.NFT) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[nft.ID] = nft
	return nil
}

func (db *DB) Delete(ctx context.Context, id uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
