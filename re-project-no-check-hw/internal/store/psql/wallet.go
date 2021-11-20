package psql

import (
	"context"
	"example.com/internal/database/psql/models"
	"example.com/internal/database/psql/store"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

func (db *DB) Wallets() store.WalletRepository {
	if db.wallets == nil {
		db.wallets = NewWalletsRepository(db.conn)
	}
	return db.wallets
}

type WalletsRepository struct {
	conn *sqlx.DB
}

func NewWalletsRepository(conn *sqlx.DB) store.WalletRepository {
	return &WalletsRepository{conn: conn}
}

func (c WalletsRepository) Create(ctx context.Context, wallet *models.Wallets) error {
	_, err := c.conn.NamedExec(`INSERT INTO Wallets(balance)
								VALUES (:balance)`, wallet)
	if err != nil {
		return err
	}
	return nil
}

func (c WalletsRepository) All(ctx context.Context) ([]*models.Wallets, error) {
	collections := make([]*models.Wallets, 0)
	if err := c.conn.Select(&collections, "SELECT * FROM Wallets"); err != nil {
		return nil, err
	}

	return collections, nil
}

func (c WalletsRepository) ByID(ctx context.Context, id int) (*models.Wallets, error) {
	collection := new(models.Wallets)
	if err := c.conn.Get(collection, "SELECT * FROM Wallets WHERE id=$1", id); err != nil {
		return nil, err
	}

	return collection, nil
}

func (c WalletsRepository) Update(ctx context.Context, collection *models.Wallets) error {
	var query []string
	v := reflect.ValueOf(*collection)
	typeOf := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == reflect.Zero(reflect.TypeOf(v.Field(i).Interface())).Interface() {
			continue
		}
		query = append(query, fmt.Sprintf("%s=%v", typeOf.Field(i).Name, v.Field(i).Interface()))
	}
	_, err := c.conn.Exec("UPDATE Wallets SET $1 WHERE id = $2", query, collection.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c WalletsRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM Wallets WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
