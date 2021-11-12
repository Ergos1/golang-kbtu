package postgres

import (
	"example.com/internal/store/psql/store"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

func getEnv() (map[string]string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("[error] Error during getting path")
	}
	godotenv.Load(pwd + "\\.env")
	env, err := godotenv.Read()
	if err != nil {
		return nil, fmt.Errorf("[error] Error during reading .env file")
	}
	return env, nil
}

func getPsqlInfo() (string, error) {
	env, err := getEnv()
	if err != nil {
		return "", err
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		env["HOST"], env["PORT"], env["USER"], env["PASSWORD"], env["DBNAME"])
	return psqlInfo, nil
}

type DB struct {
	conn *sqlx.DB

	collections       store.CollectionRepository
	nonFungibleTokens store.NonFungibleTokenRepository
	transactions      store.TransactionRepository
	users             store.UserRepository
	wallets           store.WalletRepository
}


func NewDB() store.Store {
	return &DB{}
}

func (db *DB) Connect() error {
	psqlInfo, err := getPsqlInfo()
	if err != nil {
		return err
	}
	conn, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}

	db.conn = conn
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
