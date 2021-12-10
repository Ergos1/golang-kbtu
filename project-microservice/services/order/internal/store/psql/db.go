package psql

import (
	"example.com/services/order/internal/store"
	"example.com/pkg/database/psql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	conn *sqlx.DB

	orders store.OrdersRepository
}

func NewDB() store.Store {
	return &DB{}
}

func (db *DB) Connect(uri string) error {
	conn, err := psql.NewConnection(uri)
	if err != nil {
		return err
	}

	db.conn = conn
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
