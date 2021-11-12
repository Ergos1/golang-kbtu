package postgres

import (
	"context"
	"example.com/internal/store/psql/models"
	"example.com/internal/store/psql/store"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

func (db *DB) Users() store.UserRepository {
	if db.users == nil {
		db.users = NewUsersRepository(db.conn)
	}
	return db.users
}

type UsersRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) store.UserRepository {
	return &UsersRepository{conn: conn}
}

func (c UsersRepository) Create(ctx context.Context, collection *models.Client) error {
	_, err := c.conn.NamedExec(`INSERT INTO Users(name, symbol, description, ownerid)
								VALUES (:name, :symbol, :description, :ownerid)`, collection)
	if err != nil {
		return err
	}
	return nil
}

func (c UsersRepository) All(ctx context.Context) ([]*models.Client, error) {
	collections := make([]*models.Client, 0)
	if err := c.conn.Select(&collections, "SELECT * FROM Users"); err != nil {
		return nil, err
	}

	return collections, nil
}

func (c UsersRepository) ByID(ctx context.Context, id int) (*models.Client, error) {
	collection := new(models.Client)
	if err := c.conn.Get(collection, "SELECT * FROM Users WHERE id=$1", id); err != nil {
		return nil, err
	}

	return collection, nil
}

func (c UsersRepository) Update(ctx context.Context, collection *models.Client) error {
	var query []string
	v := reflect.ValueOf(*collection)
	typeOf := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == reflect.Zero(reflect.TypeOf(v.Field(i).Interface())).Interface() {
			continue
		}
		query = append(query, fmt.Sprintf("%s=%v", typeOf.Field(i).Name, v.Field(i).Interface()))
	}
	_, err := c.conn.Exec("UPDATE Users SET $1 WHERE id = $2", query, collection.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c UsersRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM Users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
