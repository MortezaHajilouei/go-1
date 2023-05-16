package redis

import (
	"context"
	"encoding/json"
	"log"
	"micro/config"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// rds struct for redis adapter
type rds struct {
	db *redis.Client
}

func NewRedis(lc fx.Lifecycle) Store {
	rds := rds{}
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			log.Println("redis connected")
			return rds.connect(*config.C())
		},
		OnStop: func(c context.Context) error {
			log.Println("redis closed")
			return rds.db.Close()
		},
	})
	return &rds
}

// Connect, method for connect to redis
func (r *rds) connect(confs config.Config) error {
	var err error

	r.db = redis.NewClient(&redis.Options{
		DB:       confs.Redis.DB,
		Addr:     confs.Redis.Host,
		Username: confs.Redis.Username,
		Password: confs.Redis.Password,
	})

	if err = r.db.Ping(context.Background()).Err(); err != nil {
		zap.L().Error(err.Error())
	}

	return err
}

func (r *rds) SetClient(client *redis.Client) {
	r.db = client
}

// Set meth a new key,value
func (r *rds) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return r.db.Set(ctx, key, p, duration).Err()
}

// Get meth, get value with key
func (r *rds) Get(ctx context.Context, key string, dest interface{}) error {
	p, err := r.db.Get(ctx, key).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(p), &dest)
}

// Del for delete keys in redis
func (r *rds) Del(ctx context.Context, key ...string) error {
	_, err := r.db.Del(ctx, key...).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

// DelWithPattern delete all keys with a special pattern
func (r *rds) DelWithPattern(ctx context.Context, pattern string) error {
	// scan to find all keys with this pattern
	iter := r.db.Scan(ctx, 0, pattern, 0).Iterator()
	// iterate on found keys to delete
	var err error
	for iter.Next(ctx) {
		// call delete function per item
		err = r.Del(ctx, iter.Val())
	}
	// if there is no error this variable will be nil
	return err
}

func (r *rds) ListPush(ctx context.Context, ttl time.Duration, key string, values ...interface{}) error {
	if err := r.db.Get(ctx, key).Err(); err == redis.Nil {
		defer r.db.Expire(ctx, key, ttl)
	}
	return r.db.LPush(ctx, key, values...).Err()
}

func (r *rds) ListPop(ctx context.Context, key string) ([]interface{}, error) {
	res := r.db.LPop(ctx, key)
	return res.Args(), res.Err()
}

func (r *rds) ListRange(ctx context.Context, key string, from, to int) ([]string, error) {
	return r.db.LRange(ctx, key, int64(from), int64(to)).Result()
}
