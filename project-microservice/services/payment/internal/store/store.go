package store

import (
	"context"

	"example.com/services/payment/internal/models"
)

type Store interface {
	Connect(uri string) error
	Close() error

	Payments() PaymentsRepository
}

type PaymentsRepository interface {
	Create(ctx context.Context, collection *models.Payment) error
	All(ctx context.Context) ([]*models.Payment, error)
	ByID(ctx context.Context, id uint) (*models.Payment, error)
	Update(ctx context.Context, collection *models.Payment) error
	Delete(ctx context.Context, id uint) error
}
