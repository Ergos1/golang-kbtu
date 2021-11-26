package models

type Wallet struct {
	Id      uint    `db:"id" json:"id"`
	Balance float64 `db:"balance" json:"balance"`
}

type WalletFilter struct {
	Query *string `json:"query"`
}

// func (w Wallet) MarshalBinary() ([]byte, error) {
// 	return json.Marshal(w)
// }

// func (w *Wallet) MarshalBinary() ([]byte, error) {
// 	return json.Marshal(w)
// }
