package psql

import (
	"context"

	"example.com/services/order/internal/models"
	"example.com/services/order/internal/store"
	"example.com/pkg/database/psql/operations"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Orders() store.OrdersRepository {
	if db.orders == nil {
		db.orders = NewOrdersRepository(db.conn)
	}
	return db.orders
}

type OrdersRepository struct {
	conn *sqlx.DB
}

func NewOrdersRepository(conn *sqlx.DB) store.OrdersRepository {
	return &OrdersRepository{conn: conn}
}

func (c OrdersRepository) Create(ctx context.Context, order *models.Order) error {
	_, err := operations.Insert(c.conn, "Orders", order)
	if err != nil {
		return err
	}
	return nil
}

func (c OrdersRepository) All(ctx context.Context) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)

	if err := c.conn.Select(&orders, "SELECT * FROM Orders"); err != nil {
		return nil, err
	}
	return orders, nil
}

func (c OrdersRepository) ByID(ctx context.Context, id uint) (*models.Order, error) {
	order := new(models.Order)
	if err := c.conn.Get(order, "SELECT * FROM Orders WHERE id=$1", id); err != nil {
		return nil, err
	}

	return order, nil
}

func (c OrdersRepository) Update(ctx context.Context, order *models.Order) error {
	if _, err := operations.Update(c.conn, "Orders", order); err != nil {
		return err
	}
	return nil
}

func (c OrdersRepository) Delete(ctx context.Context, id uint) error {
	if _, err := operations.Delete(c.conn, "Orders", id); err != nil {
		return err
	}
	return nil
}
