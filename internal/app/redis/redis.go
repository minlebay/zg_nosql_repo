package redis

import (
	"context"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"sync"
)

type Redis struct {
	Config *Config
	Logger *zap.Logger
	db     *redis.Client
	wg     sync.WaitGroup
}

func NewRedis(logger *zap.Logger, config *Config) *Redis {
	return &Redis{
		Config: config,
		Logger: logger,
	}
}

func (r *Redis) StartRedis(ctx context.Context) {
	go func() {
		numdb, err := strconv.ParseInt(r.Config.DB, 10, 64)
		if err != nil {
			r.Logger.Error("Failed to parse DB", zap.Error(err))
		}

		r.db = redis.NewClient(&redis.Options{
			Addr: r.Config.Address,
			DB:   int(numdb),
		})
	}()
}

func (r *Redis) StopRedis(ctx context.Context) {
	r.wg.Wait()
	err := r.db.Close()
	if err != nil {
		r.Logger.Error("Failed to disconnect from Redis", zap.Error(err))
	}
}

func (s *Redis) Get(key string) (out []byte, err error) {
	out, err = s.db.Get(key).Bytes()
	return
}

func (s *Redis) Put(key string, value []byte) (err error) {
	err = s.db.Set(key, string(value), 0).Err()
	return
}

func (s *Redis) Delete(key string) (err error) {
	err = s.db.Del(key).Err()
	return
}

func (s *Redis) Iterate(filter string) (out []string, err error) {
	if filter != "" {
		filter += "*"
	}
	iter := s.db.Scan(0, filter, 0).Iterator()
	for iter.Next() {
		key := iter.Val()
		out = append(out, key)
	}
	sort.Strings(out)
	err = iter.Err()
	return
}
