package store

import (
	"context"

	"example.com/internal/models"
)

type Store interface {
	Connect(uri string) error
	Close() error

	Collections() CollectionsRepository
	Assets() AssetsRepository
	Transactions() TransactionsRepository
	Accounts() AccountsRepository
	Wallets() WalletsRepository
	Categories() CategoriesRepository
}

type CollectionsRepository interface {
	Create(ctx context.Context, collection *models.Collection) error
	All(ctx context.Context, filter *models.CollectionFilter) ([]*models.Collection, error)
	ByID(ctx context.Context, id uint) (*models.Collection, error)
	Update(ctx context.Context, collection *models.Collection) error
	Delete(ctx context.Context, id uint) error
}

type AssetsRepository interface {
	Create(ctx context.Context, nft *models.Asset) error
	All(ctx context.Context, filter *models.AssetFilter) ([]*models.Asset, error)
	ByID(ctx context.Context, id uint) (*models.Asset, error)
	Update(ctx context.Context, nft *models.Asset) error
	Delete(ctx context.Context, id uint) error
}

type TransactionsRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	All(ctx context.Context, filter *models.TransactionFilter) ([]*models.Transaction, error)
	ByID(ctx context.Context, id uint) (*models.Transaction, error)
	Update(ctx context.Context, transaction *models.Transaction) error
	Delete(ctx context.Context, id uint) error
}

type AccountsRepository interface {
	Create(ctx context.Context, user *models.Account) error
	All(ctx context.Context, filter *models.AccountFilter) ([]*models.Account, error)
	ByID(ctx context.Context, id uint) (*models.Account, error)
	Update(ctx context.Context, user *models.Account) error
	Delete(ctx context.Context, id uint) error
}

type WalletsRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	All(ctx context.Context, filter *models.WalletFilter) ([]*models.Wallet, error)
	ByID(ctx context.Context, id uint) (*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet) error
	Delete(ctx context.Context, id uint) error
}

type CategoriesRepository interface {
	Create(ctx context.Context, wallet *models.Category) error
	All(ctx context.Context, filter *models.CategoryFilter) ([]*models.Category, error)
	ByID(ctx context.Context, id uint) (*models.Category, error)
	Update(ctx context.Context, wallet *models.Category) error
	Delete(ctx context.Context, id uint) error
}
