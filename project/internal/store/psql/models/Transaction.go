package models

type Transaction struct {
	Id          uint    `db:"id" json:"id"`
	WalletId    uint    `db:"walletid" json:"wallet_id"`
	ToUserId    uint    `db:"touserid" json:"to_user_id"`
	Amount      float64 `db:"amount" json:"amount"`
	Description string  `db:"description" json:"description"`
}
