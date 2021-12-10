package psql

import (
	"context"

	"example.com/internal/models"
	"example.com/internal/store"
	"example.com/pkg/database/psql/operations"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Accounts() store.AccountsRepository {
	if db.accounts == nil {
		db.accounts = NewAccountsRepository(db.conn)
	}
	return db.accounts
}

type AccountsRepository struct {
	conn *sqlx.DB
}

func NewAccountsRepository(conn *sqlx.DB) store.AccountsRepository {
	return &AccountsRepository{conn: conn}
}

func (c AccountsRepository) Create(ctx context.Context, account *models.Account) error {
	_, err := operations.Insert(c.conn, "Accounts", account)
	if err != nil {
		return err
	}
	return nil
}

func (c AccountsRepository) All(ctx context.Context, filter *models.AccountFilter) ([]*models.Account, error) {
	accounts := make([]*models.Account, 0)
	if filter.Username != nil {
		if err := c.conn.Select(&accounts, "SELECT * FROM Account WHERE username ILIKE $1", "%"+*filter.Username+"%"); err != nil {
			return nil, err
		}

		return accounts, nil
	}
	if err := c.conn.Select(&accounts, "SELECT * FROM Accounts"); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (c AccountsRepository) ByID(ctx context.Context, id uint) (*models.Account, error) {
	account := new(models.Account)
	if err := c.conn.Get(account, "SELECT * FROM Accounts WHERE id=$1", id); err != nil {
		return nil, err
	}

	return account, nil
}

func (c AccountsRepository) Update(ctx context.Context, account *models.Account) error {
	if _, err := operations.Update(c.conn, "Accounts", account); err != nil {
		return err
	}
	return nil
}

func (c AccountsRepository) Delete(ctx context.Context, id uint) error {
	if _, err := operations.Delete(c.conn, "Accounts", id); err != nil {
		return err
	}
	return nil
}
