package store

import (
	"context"

	"example.com/internal/models"
)

type Store interface {
	Create(ctx context.Context, nft *models.NFT) error
	All(ctx context.Context) ([]*models.NFT, error)
	ByID(ctx context.Context, id uint64) (*models.NFT, error)
	Update(ctx context.Context, NFT *models.NFT) error
	Delete(ctx context.Context, id uint64) error
}
