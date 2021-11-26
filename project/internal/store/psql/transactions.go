package psql


import (
	"context"

	"example.com/internal/models"
	"example.com/internal/store"
	"example.com/pkg/database/psql/operations"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Transactions() store.TransactionsRepository {
	if db.transactions == nil {
		db.transactions = NewTransactionsRepository(db.conn)
	}
	return db.transactions
}

type TransactionsRepository struct {
	conn *sqlx.DB
}

func NewTransactionsRepository(conn *sqlx.DB) store.TransactionsRepository {
	return &TransactionsRepository{conn: conn}
}

func (c TransactionsRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	_, err := operations.Insert(c.conn, "Transactions", transaction)
	if err != nil {
		return err
	}
	return nil
}

func (c TransactionsRepository) All(ctx context.Context, filter *models.TransactionFilter) ([]*models.Transaction, error) {
	transactions := make([]*models.Transaction, 0)
	if filter.Query != nil {
		if err := c.conn.Select(&transactions, "SELECT * FROM Transactions WHERE Description ILIKE $1", "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return transactions, nil
	}
	if err := c.conn.Select(&transactions, "SELECT * FROM Transactions"); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (c TransactionsRepository) ByID(ctx context.Context, id uint) (*models.Transaction, error) {
	transaction := new(models.Transaction)
	if err := c.conn.Get(transaction, "SELECT * FROM Transactions WHERE id=$1", id); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (c TransactionsRepository) Update(ctx context.Context, transaction *models.Transaction) error {
	if _, err := operations.Update(c.conn, "Transactions", transaction); err != nil {
		return err
	}
	return nil
}

func (c TransactionsRepository) Delete(ctx context.Context, id uint) error {
	if _, err := operations.Delete(c.conn, "Transactions", id); err != nil {
		return err
	}
	return nil
}
