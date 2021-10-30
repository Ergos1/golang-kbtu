package models

type Client struct {
	Id       uint	`db:"id"`
	WalletId uint	`db:"walletid"`
	Username string	`db:"username"`
	Email    string	`db:"email"`
	Password string	`db:"password"`
}
