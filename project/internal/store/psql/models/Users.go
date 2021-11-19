package models

type Clients struct {
	Id       uint	`db:"id" json:"id"`
	WalletId uint	`db:"walletid" json:"wallet_id"`
	Username string	`db:"username" json:"username"`
	Email    string	`db:"email" json:"email"`
	Password string	`db:"password" json:"password"`
}
