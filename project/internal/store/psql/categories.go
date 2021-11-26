package psql


import (
	"context"

	"example.com/internal/models"
	"example.com/internal/store"
	"example.com/pkg/database/psql/operations"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Categories() store.CategoriesRepository {
	if db.categories == nil {
		db.categories = NewCategoriesRepository(db.conn)
	}
	return db.categories
}

type CategoriesRepository struct {
	conn *sqlx.DB
}

func NewCategoriesRepository(conn *sqlx.DB) store.CategoriesRepository {
	return &CategoriesRepository{conn: conn}
}

func (c CategoriesRepository) Create(ctx context.Context, category *models.Category) error {
	_, err := operations.Insert(c.conn, "Categories", category)
	if err != nil {
		return err
	}
	return nil
}

func (c CategoriesRepository) All(ctx context.Context, filter *models.CategoryFilter) ([]*models.Category, error) {
	categories := make([]*models.Category, 0)
	if filter.Query != nil {
		if err := c.conn.Select(&categories, "SELECT * FROM Categories WHERE name ILIKE $1", "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return categories, nil
	}
	if err := c.conn.Select(&categories, "SELECT * FROM Categories"); err != nil {
		return nil, err
	}
	return categories, nil
}

func (c CategoriesRepository) ByID(ctx context.Context, id uint) (*models.Category, error) {
	category := new(models.Category)
	if err := c.conn.Get(category, "SELECT * FROM Categories WHERE id=$1", id); err != nil {
		return nil, err
	}

	return category, nil
}

func (c CategoriesRepository) Update(ctx context.Context, category *models.Category) error {
	if _, err := operations.Update(c.conn, "Categories", category); err != nil {
		return err
	}
	return nil
}

func (c CategoriesRepository) Delete(ctx context.Context, id uint) error {
	if _, err := operations.Delete(c.conn, "Categories", id); err != nil {
		return err
	}
	return nil
}
