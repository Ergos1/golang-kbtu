package cache

import (
	"context"
	"time"
)

type Cache interface {
	Connect(host string, db int, expires time.Duration)
	Close() error

	Get(ctx context.Context, key interface{}, dest interface{}) error
	Set(ctx context.Context, key interface{}, value interface{}) error
	Remove(ctx context.Context, key interface{}) error
	Purge(ctx context.Context) error
}
