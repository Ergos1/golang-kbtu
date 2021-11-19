package postgres

import (
	"context"
	"example.com/internal/store/psql/models"
	"example.com/internal/store/psql/store"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

func (db *DB) NonFungibleTokens() store.NonFungibleTokenRepository {
	if db.nonFungibleTokens == nil {
		db.nonFungibleTokens = NewNonFungibleTokensRepository(db.conn)
	}
	return db.nonFungibleTokens
}

type NonFungibleTokensRepository struct {
	conn *sqlx.DB
}

func NewNonFungibleTokensRepository(conn *sqlx.DB) store.NonFungibleTokenRepository {
	return &NonFungibleTokensRepository{conn: conn}
}

func (c NonFungibleTokensRepository) Create(ctx context.Context, collection *models.NonFungibleTokens) error {
	_, err := c.conn.NamedExec(`INSERT INTO NonFungibleTokens(name, symbol, description, ownerid)
								VALUES (:name, :symbol, :description, :ownerid)`, collection)
	if err != nil {
		return err
	}
	return nil
}

func (c NonFungibleTokensRepository) All(ctx context.Context) ([]*models.NonFungibleTokens, error) {
	collections := make([]*models.NonFungibleTokens, 0)
	if err := c.conn.Select(&collections, "SELECT * FROM NonFungibleTokens"); err != nil {
		return nil, err
	}

	return collections, nil
}

func (c NonFungibleTokensRepository) ByID(ctx context.Context, id int) (*models.NonFungibleTokens, error) {
	collection := new(models.NonFungibleTokens)
	if err := c.conn.Get(collection, "SELECT * FROM NonFungibleTokens WHERE id=$1", id); err != nil {
		return nil, err
	}

	return collection, nil
}

func (c NonFungibleTokensRepository) Update(ctx context.Context, collection *models.NonFungibleTokens) error {
	var query []string
	v := reflect.ValueOf(*collection)
	typeOf := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == reflect.Zero(reflect.TypeOf(v.Field(i).Interface())).Interface() {
			continue
		}
		query = append(query, fmt.Sprintf("%s=%v", typeOf.Field(i).Name, v.Field(i).Interface()))
	}
	_, err := c.conn.Exec("UPDATE NonFungibleTokens SET $1 WHERE id = $2", query, collection.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c NonFungibleTokensRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM NonFungibleTokens WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
