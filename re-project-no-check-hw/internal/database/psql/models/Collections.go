package models

type Collections struct {
	Id          uint	`db:"id" json:"id"`
	Name        string	`db:"name" json:"name"`
	Symbol      string	`db:"symbol" json:"symbol"`
	Description string	`db:"description" json:"description"`
	OwnerId     uint	`db:"ownerid" json:"owner_id"`
}
