package redis

import "example.com/internal/database/psql/models"

type UsersCache interface {
	Set(key string, value *models.Clients)
	Get(key string) *models.Clients
	SetAll(key string, values []*models.Clients)
	GetAll(key string) []*models.Clients
	Delete(key string)
}