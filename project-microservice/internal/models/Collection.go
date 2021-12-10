package models

type Collection struct {
	Id          uint    `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	CreatorId   uint    `db:"creator_id" json:"creator_id"`
	// Assets      []Asset `db:"assets" json:"assets"` one2many
} 

type CollectionFilter struct {
	Query *string `json:"query"`
}
