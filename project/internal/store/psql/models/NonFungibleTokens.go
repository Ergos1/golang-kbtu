package models

type NonFungibleTokens struct {
	Id           uint    `db:"id" json:"id"`
	Likes        uint    `db:"likes" json:"likes"`
	CollectionId uint    `db:"collectionid" json:"collection_id"`
	OwnerId      uint    `db:"ownerid" json:"owner_id"`
	Price        float64 `db:"price" json:"price"`
	Royalties    uint    `db:"royalties" json:"royalties"`
	Title        string  `db:"title" json:"title"`
	Description  string  `db:"description" json:"description"`
}
