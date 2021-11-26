package psql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewConnection(uri string) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("[error] Error occurred while connecting to postgres")
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}

