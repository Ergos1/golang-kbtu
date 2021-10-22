package models

type NFT struct {
	ID       uint64            `json:"id"`
	Name     string            `json:"name"`
	Symbol   string            `json:"symbol"`
	TokenURI map[uint64]string `json:"tokenUri"`
}
