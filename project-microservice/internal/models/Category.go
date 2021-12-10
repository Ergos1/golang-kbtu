package models

type Category struct {
	Id     uint    `db:"id" json:"id"`
	Name   string  `db:"name" json:"name"`
}

type CategoryFilter struct {
	Query *string `json:"query"`
}
