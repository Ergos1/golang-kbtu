package models

type Transaction struct {
	Id           uint    `db:"id" json:"id"`
	FromWalletId uint    `db:"from_wallet_id" json:"from_wallet_id"`
	ToWalletId   uint    `db:"to_wallet_id" json:"to_wallet_id"`
	Amount       float64 `db:"amount" json:"amount"`
	Description  string  `db:"description" json:"description"`
}

type TransactionFilter struct {
	Query *string `json:"query"`
}

