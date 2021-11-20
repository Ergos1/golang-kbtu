package redis

import "example.com/internal/database/psql/models"

type WalletsCache interface {
	Set(key string, value *models.Wallets)
	Get(key string) *models.Wallets
	Delete(key string)
}