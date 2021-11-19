package postgres

import (
	"context"
	"example.com/internal/store/psql/models"
	"example.com/internal/store/psql/store"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

func(db *DB) Collections() store.CollectionRepository {
	if db.collections == nil {
		db.collections = NewCollectionsRepository(db.conn)
	}

	return db.collections
}

type CollectionsRepository struct {
	conn *sqlx.DB
}

func NewCollectionsRepository(conn *sqlx.DB) store.CollectionRepository {
	return &CollectionsRepository{conn: conn}
}

func (c CollectionsRepository) Create(ctx context.Context, collection *models.Collections) error {
	_, err := c.conn.NamedExec(`INSERT INTO Collections(name, symbol, description, ownerid)
								VALUES (:name, :symbol, :description, :ownerid)`, collection)
	if err != nil {
		return err
	}
	return nil
}

func (c CollectionsRepository) All(ctx context.Context) ([]*models.Collections, error) {
	collections := make([]*models.Collections, 0)
	if err:=c.conn.Select(&collections, "SELECT * FROM collections"); err != nil {
		return nil, err
	}

	return collections, nil
}

func (c CollectionsRepository) ByID(ctx context.Context, id int) (*models.Collections, error) {
	collection := new(models.Collections)
	if err := c.conn.Get(collection, "SELECT * FROM categories WHERE id=$1", id); err != nil {
		return nil, err
	}

	return collection, nil
}

func (c CollectionsRepository) Update(ctx context.Context, collection *models.Collections) error {
	var query []string
	v := reflect.ValueOf(*collection)
	typeOf := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == reflect.Zero(reflect.TypeOf(v.Field(i).Interface())).Interface(){
			continue
		}
		query = append(query, fmt.Sprintf("%s=%v", typeOf.Field(i).Name, v.Field(i).Interface()))
	}
	_, err := c.conn.Exec("UPDATE collections SET $1 WHERE id = $2", query, collection.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c CollectionsRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM collections WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}


