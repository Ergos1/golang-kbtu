package postgres

import (
	"context"
	"example.com/internal/store/psql/models"
	"example.com/internal/store/psql/store"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

func (db *DB) Transactions() store.TransactionRepository {
	if db.transactions == nil {
		db.transactions = NewTransactionsRepository(db.conn)
	}
	return db.transactions
}

type TransactionsRepository struct {
	conn *sqlx.DB
}

func NewTransactionsRepository(conn *sqlx.DB) store.TransactionRepository {
	return &TransactionsRepository{conn: conn}
}

func (c TransactionsRepository) Create(ctx context.Context, collection *models.Transactions) error {
	_, err := c.conn.NamedExec(`INSERT INTO Transactions(name, symbol, description, ownerid)
								VALUES (:name, :symbol, :description, :ownerid)`, collection)
	if err != nil {
		return err
	}
	return nil
}

func (c TransactionsRepository) All(ctx context.Context) ([]*models.Transactions, error) {
	collections := make([]*models.Transactions, 0)
	if err := c.conn.Select(&collections, "SELECT * FROM Transactions"); err != nil {
		return nil, err
	}

	return collections, nil
}

func (c TransactionsRepository) ByID(ctx context.Context, id int) (*models.Transactions, error) {
	collection := new(models.Transactions)
	if err := c.conn.Get(collection, "SELECT * FROM Transactions WHERE id=$1", id); err != nil {
		return nil, err
	}

	return collection, nil
}

func (c TransactionsRepository) Update(ctx context.Context, collection *models.Transactions) error {
	var query []string
	v := reflect.ValueOf(*collection)
	typeOf := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == reflect.Zero(reflect.TypeOf(v.Field(i).Interface())).Interface() {
			continue
		}
		query = append(query, fmt.Sprintf("%s=%v", typeOf.Field(i).Name, v.Field(i).Interface()))
	}
	_, err := c.conn.Exec("UPDATE Transactions SET $1 WHERE id = $2", query, collection.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c TransactionsRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM Transactions WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
