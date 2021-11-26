package psql

import (
	"context"

	"example.com/internal/models"
	"example.com/internal/store"
	"example.com/pkg/database/psql/operations"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Assets() store.AssetsRepository {
	if db.assets == nil {
		db.assets = NewAssetsRepository(db.conn)
	}
	return db.assets
}

type AssetsRepository struct {
	conn *sqlx.DB
}

func NewAssetsRepository(conn *sqlx.DB) store.AssetsRepository {
	return &AssetsRepository{conn: conn}
}

func (c AssetsRepository) Create(ctx context.Context, asset *models.Asset) error {
	_, err := operations.Insert(c.conn, "Assets", asset)
	if err != nil {
		return err
	}
	return nil
}

func (c AssetsRepository) All(ctx context.Context, filter *models.AssetFilter) ([]*models.Asset, error) {
	assets := make([]*models.Asset, 0)
	if filter.Query != nil {
		if err := c.conn.Select(&assets, "SELECT * FROM Assets WHERE title ILIKE $1", "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return assets, nil
	}
	if err := c.conn.Select(&assets, "SELECT * FROM Assets", filter.Query); err != nil {
		return nil, err
	}
	return assets, nil
}

func (c AssetsRepository) ByID(ctx context.Context, id uint) (*models.Asset, error) {
	asset := new(models.Asset)
	if err := c.conn.Get(asset, "SELECT * FROM Assets WHERE id=$1", id); err != nil {
		return nil, err
	}
	return asset, nil
}

func (c AssetsRepository) Update(ctx context.Context, asset *models.Asset) error {
	if _, err := operations.Update(c.conn, "Assets", asset); err != nil {
		return err
	}
	return nil
}

func (c AssetsRepository) Delete(ctx context.Context, id uint) error {
	if _, err := operations.Delete(c.conn, "Assets", id); err != nil {
		return err
	}
	return nil
}
