package models

type NonFungibleToken struct {
	Id           uint    `db:"id"`
	Likes        uint    `db:"likes"`
	CollectionId uint    `db:"collectionid"`
	OwnerId      uint    `db:"ownerid"`
	Price        float64 `db:"price"`
	Royalties    uint    `db:"royalties"`
	Title        string  `db:"title"`
	Description  string  `db:"description"`
}
