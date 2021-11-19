package cache

import "example.com/internal/store/psql/models"

type WalletsCache interface {
	Set(key string, value *models.Wallets)
	Get(key string) *models.Wallets
	Delete(key string)
}