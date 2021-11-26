package models

type Account struct {
	Id       uint   `db:"id" json:"id"`
	WalletId uint   `db:"wallet_id" json:"wallet_id"` // unique
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type AccountFilter struct {
	Query *string `json:"query"`
}
