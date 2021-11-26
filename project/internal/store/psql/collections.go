package psql

import (
	"context"

	"example.com/internal/models"
	"example.com/internal/store"
	"example.com/pkg/database/psql/operations"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Collections() store.CollectionsRepository {
	if db.collections == nil {
		db.collections = NewCollectionsRepository(db.conn)
	}
	return db.collections
}

type CollectionsRepository struct {
	conn *sqlx.DB
}

func NewCollectionsRepository(conn *sqlx.DB) store.CollectionsRepository {
	return &CollectionsRepository{conn: conn}
}

func (c CollectionsRepository) Create(ctx context.Context, collection *models.Collection) error {
	_, err := operations.Insert(c.conn, "Collections", collection)
	if err != nil {
		return err
	}
	return nil
}

func (c CollectionsRepository) All(ctx context.Context, filter *models.CollectionFilter) ([]*models.Collection, error) {
	collections := make([]*models.Collection, 0)
	if filter.Query != nil {
		if err := c.conn.Select(&collections, "SELECT * FROM Collections WHERE name ILIKE $1", "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return collections, nil
	}
	if err := c.conn.Select(&collections, "SELECT * FROM Collections"); err != nil {
		return nil, err
	}
	return collections, nil
}

func (c CollectionsRepository) ByID(ctx context.Context, id uint) (*models.Collection, error) {
	collection := new(models.Collection)
	if err := c.conn.Get(collection, "SELECT * FROM Collections WHERE id=$1", id); err != nil {
		return nil, err
	}

	return collection, nil
}

func (c CollectionsRepository) Update(ctx context.Context, collection *models.Collection) error {
	if _, err := operations.Update(c.conn, "Collections", collection); err != nil {
		return err
	}
	return nil
}

func (c CollectionsRepository) Delete(ctx context.Context, id uint) error {
	if _, err := operations.Delete(c.conn, "Collections", id); err != nil {
		return err
	}
	return nil
}
