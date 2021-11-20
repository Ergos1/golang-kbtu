package models

type Wallets struct {
	Id      uint    `db:"id" json:"id"`
	Balance float64 `db:"balance" json:"balance"`
}
