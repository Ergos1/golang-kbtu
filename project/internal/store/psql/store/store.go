package store

import (
	"context"
	"example.com/internal/store/psql/models"
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
	Create(ctx context.Context, collection *models.Collection) error
	All(ctx context.Context) ([]*models.Collection, error)
	ByID(ctx context.Context, id int) (*models.Collection, error)
	Update(ctx context.Context, collection *models.Collection) error
	Delete(ctx context.Context, id int) error
}

type NonFungibleTokenRepository interface {
	Create(ctx context.Context, nft *models.NonFungibleToken) error
	All(ctx context.Context) ([]*models.NonFungibleToken, error)
	ByID(ctx context.Context, id int) (*models.NonFungibleToken, error)
	Update(ctx context.Context, nft *models.NonFungibleToken) error
	Delete(ctx context.Context, id int) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	All(ctx context.Context) ([]*models.Transaction, error)
	ByID(ctx context.Context, id int) (*models.Transaction, error)
	Update(ctx context.Context, transaction *models.Transaction) error
	Delete(ctx context.Context, id int) error
}

type UserRepository interface {
	Create(ctx context.Context, user *models.Client) error
	All(ctx context.Context) ([]*models.Client, error)
	ByID(ctx context.Context, id int) (*models.Client, error)
	Update(ctx context.Context, user *models.Client) error
	Delete(ctx context.Context, id int) error
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	All(ctx context.Context) ([]*models.Wallet, error)
	ByID(ctx context.Context, id int) (*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet) error
	Delete(ctx context.Context, id int) error
}


