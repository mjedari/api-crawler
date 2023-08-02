package contracts

import (
	"context"
	"time"
)

type IStorage interface {
	Store(ctx context.Context, key, value string, timeToLive time.Duration) error
	BatchStore(ctx context.Context, data map[string]string, timeToLive time.Duration) error
	Fetch(ctx context.Context, key string) []byte
	FetchAll(ctx context.Context, pattern string) [][]byte
	Exists(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
}
