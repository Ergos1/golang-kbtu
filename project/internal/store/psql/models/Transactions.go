package models

type Transaction struct {
	Id          uint    `db:"id"`
	WalletId    uint    `db:"walletid"`
	ToUserId    uint    `db:"touserid"`
	Amount      float64 `db:"amount"`
	Description string  `db:"description"`
}
