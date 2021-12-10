package store

import (
	"context"

	"example.com/services/order/internal/models"
)

type Store interface {
	Connect(uri string) error
	Close() error

	Orders() OrdersRepository
}

type OrdersRepository interface {
	Create(ctx context.Context, collection *models.Order) error
	All(ctx context.Context) ([]*models.Order, error)
	ByID(ctx context.Context, id uint) (*models.Order, error)
	Update(ctx context.Context, collection *models.Order) error
	Delete(ctx context.Context, id uint) error
}
