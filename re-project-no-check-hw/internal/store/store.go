package store

import (
	"context"
	"example.com/internal/database/psql/models"
)

type Store interface {
	Connect() error
	Close() error

	Collections() CollectionRepository
	NonFungibleTokens() NonFungibleTokenRepository
	Transactions() TransactionRepository
	Users() UserRepository
	Wallets() WalletRepository
}

type CollectionRepository interface {
	Create(ctx context.Context, collection *models.Collections) error
	All(ctx context.Context) ([]*models.Collections, error)
	ByID(ctx context.Context, id int) (*models.Collections, error)
	Update(ctx context.Context, collection *models.Collections) error
	Delete(ctx context.Context, id int) error
}

type NonFungibleTokenRepository interface {
	Create(ctx context.Context, nft *models.NonFungibleTokens) error
	All(ctx context.Context) ([]*models.NonFungibleTokens, error)
	ByID(ctx context.Context, id int) (*models.NonFungibleTokens, error)
	Update(ctx context.Context, nft *models.NonFungibleTokens) error
	Delete(ctx context.Context, id int) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transactions) error
	All(ctx context.Context) ([]*models.Transactions, error)
	ByID(ctx context.Context, id int) (*models.Transactions, error)
	Update(ctx context.Context, transaction *models.Transactions) error
	Delete(ctx context.Context, id int) error
}

type UserRepository interface {
	Create(ctx context.Context, user *models.Clients) error
	All(ctx context.Context) ([]*models.Clients, error)
	ByID(ctx context.Context, id int) (*models.Clients, error)
	Update(ctx context.Context, user *models.Clients) error
	Delete(ctx context.Context, id int) error
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *models.Wallets) error
	All(ctx context.Context) ([]*models.Wallets, error)
	ByID(ctx context.Context, id int) (*models.Wallets, error)
	Update(ctx context.Context, wallet *models.Wallets) error
	Delete(ctx context.Context, id int) error
}


