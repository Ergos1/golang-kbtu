package models

type Wallet struct {
	Id      uint    `db:"id"`
	Balance float64 `db:"balance"`
}
