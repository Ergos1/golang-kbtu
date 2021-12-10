package psql

import (
	"context"

	"example.com/services/payment/internal/models"
	"example.com/services/payment/internal/store"
	"example.com/pkg/database/psql/operations"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Payments() store.PaymentsRepository {
	if db.payments == nil {
		db.payments = NewPaymentsRepository(db.conn)
	}
	return db.payments
}

type PaymentsRepository struct {
	conn *sqlx.DB
}

func NewPaymentsRepository(conn *sqlx.DB) store.PaymentsRepository {
	return &PaymentsRepository{conn: conn}
}

func (c PaymentsRepository) Create(ctx context.Context, payment *models.Payment) error {
	_, err := operations.Insert(c.conn, "Payments", payment)
	if err != nil {
		return err
	}
	return nil
}

func (c PaymentsRepository) All(ctx context.Context) ([]*models.Payment, error) {
	payments := make([]*models.Payment, 0)

	if err := c.conn.Select(&payments, "SELECT * FROM Payments"); err != nil {
		return nil, err
	}
	return payments, nil
}

func (c PaymentsRepository) ByID(ctx context.Context, id uint) (*models.Payment, error) {
	payment := new(models.Payment)
	if err := c.conn.Get(payment, "SELECT * FROM Payments WHERE id=$1", id); err != nil {
		return nil, err
	}

	return payment, nil
}

func (c PaymentsRepository) Update(ctx context.Context, payment *models.Payment) error {
	if _, err := operations.Update(c.conn, "Payments", payment); err != nil {
		return err
	}
	return nil
}

func (c PaymentsRepository) Delete(ctx context.Context, id uint) error {
	if _, err := operations.Delete(c.conn, "Payments", id); err != nil {
		return err
	}
	return nil
}
