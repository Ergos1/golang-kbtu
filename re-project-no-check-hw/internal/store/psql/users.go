package psql

import (
	"context"
	"example.com/internal/database/psql/models"
	"example.com/internal/database/psql/store"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
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

func (c UsersRepository) Create(ctx context.Context, collection *models.Clients) error {
	_, err := c.conn.NamedExec(`INSERT INTO CLIENTS(username, email, password, walletid)
								VALUES (:username, :email, :password, :walletid)`, collection)
	if err != nil {
		return err
	}
	return nil
}

func (c UsersRepository) All(ctx context.Context) ([]*models.Clients, error) {
	collections := make([]*models.Clients, 0)
	if err := c.conn.Select(&collections, "SELECT * FROM CLIENTS"); err != nil {
		return nil, err
	}

	return collections, nil
}

func (c UsersRepository) ByID(ctx context.Context, id int) (*models.Clients, error) {
	collection := new(models.Clients)
	if err := c.conn.Get(collection, "SELECT * FROM CLIENTS WHERE id=$1", id); err != nil {
		return nil, err
	}

	return collection, nil
}

func (c UsersRepository) Update(ctx context.Context, collection *models.Clients) error {
	var query string
	v := reflect.ValueOf(*collection)
	typeOf := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == reflect.Zero(reflect.TypeOf(v.Field(i).Interface())).Interface() {
			continue
		}
		query += fmt.Sprintf("%s=%v", typeOf.Field(i).Name, v.Field(i).Interface())
		query += ","
	}
	query = query[:len(query)-1]
	log.Println(query)
	_, err := c.conn.Exec("UPDATE CLIENTS SET $1 WHERE id = $2", query, collection.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c UsersRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM CLIENTS WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
