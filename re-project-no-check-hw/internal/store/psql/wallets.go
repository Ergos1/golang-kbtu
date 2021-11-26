package psql

import (
	"context"
	"fmt"

	"example.com/internal/models"
	"example.com/internal/store"
	"example.com/pkg/database/psql/operations"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Wallets() store.WalletsRepository {
	if db.wallets == nil {
		db.wallets = NewWalletsRepository(db.conn)
	}
	return db.wallets
}

type WalletsRepository struct {
	conn *sqlx.DB
}

func NewWalletsRepository(conn *sqlx.DB) store.WalletsRepository {
	return &WalletsRepository{conn: conn}
}

func (c WalletsRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	_, err := operations.Insert(c.conn, "Wallets", wallet)
	if err != nil {
		return err
	}
	return nil
}

func (c WalletsRepository) All(ctx context.Context, filter *models.WalletFilter) ([]*models.Wallet, error) {
	wallets := make([]*models.Wallet, 0)
	if filter.Query != nil {
		if err := c.conn.Select(&wallets, "SELECT * FROM Wallets WHERE Balance = $1", *filter.Query); err != nil {
			return nil, err
		}

		return wallets, nil
	}
	if err := c.conn.Select(&wallets, "SELECT * FROM Wallets"); err != nil {
		return nil, err
	}
	return wallets, nil
}

func (c WalletsRepository) ByID(ctx context.Context, id uint) (*models.Wallet, error) {
	wallet := new(models.Wallet)
	if err := c.conn.Get(wallet, "SELECT * FROM Wallets WHERE id=$1", id); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (c WalletsRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	if _, err := c.ByID(ctx, wallet.Id); err != nil {
		return fmt.Errorf("[error] This wallet does not exist")
	}
	if _, err := operations.Update(c.conn, "Wallets", wallet); err != nil {
		return err
	}
	return nil
}

func (c WalletsRepository) Delete(ctx context.Context, id uint) error {
	fmt.Println(id)
	if _, err := operations.Delete(c.conn, "Wallets", id); err != nil {
		return err
	}
	return nil
}
