package models

type Asset struct {
	Id           uint       `db:"id" json:"id"`
	CreatorId    uint       `db:"creator_id" json:"creator_id"`
	OwnerId      uint       `db:"owner_id" json:"owner_id"`
	CollectionId uint       `db:"collection_id" json:"collection_id"`
	Title        string     `db:"title" json:"title"`
	Description  string     `db:"description" json:"description"`
	Price        float64    `db:"price" json:"price"`
	Royalties    uint       `db:"royalties" json:"royalties"`
	// Categories   []Category `db:"categories" json:"categories"` many2many
}

type AssetFilter struct {
	Query *string `json:"query"`
}

