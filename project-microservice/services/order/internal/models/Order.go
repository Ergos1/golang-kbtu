package models

type Order struct {
	Id        uint   `db:"id" json:"id"`
	AssetId   uint   `db:"asset_id" json:"asset_id"`
	AccountId uint   `db:"account_id" json:"account_id"`
	Status    Status `db:"status" json:"status"`
}
