package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// store interface is interface for store things into redis
type Store interface {
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key ...string) error
	SetClient(client *redis.Client)
	ListPush(context.Context, time.Duration, string, ...interface{}) error
	ListPop(ctx context.Context, key string) ([]interface{}, error)
	ListRange(ctx context.Context, key string, from, to int) ([]string, error)
	DelWithPattern(ctx context.Context, pattern string) error
}

type AuthStore interface {
	IsExist(ctx context.Context, key string) (bool, error)
}
