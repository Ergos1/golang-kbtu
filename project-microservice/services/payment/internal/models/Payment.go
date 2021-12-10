package models

type Payment struct {
	Id        uint   `db:"id" json:"id"`
	AssetId   uint   `db:"asset_id" json:"asset_id"`
	AccountId uint   `db:"account_id" json:"account_id"`
	Amount    float64 `db:"amount" json:"amount"`
}
