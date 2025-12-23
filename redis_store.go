package telegram

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const message_cache_key = "telegram:message"

type RedisStore struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewRedisStore(host, passWord string, dbNum int) *RedisStore {
	cli := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: passWord,
		DB:       dbNum,
	})
	ctx := context.Background()
	cmd := cli.Ping(ctx)
	if cmd.Err() != nil {
		botLog.Printf("redis connect error: %v", cmd.Err())
		return nil
	}

	return &RedisStore{
		redisClient: cli,
		ctx:         context.Background(),
	}
}

func (rs *RedisStore) RPush(value string) error {
	return rs.redisClient.RPush(rs.ctx, message_cache_key, value).Err()
}
func (rs *RedisStore) BLPop() (string, error) {
	cmd := rs.redisClient.BLPop(rs.ctx, 0, message_cache_key)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Val()[1], cmd.Err()
}
func (rs *RedisStore) Close() error {
	return rs.redisClient.Close()
}
func (rs *RedisStore) Size() int64 {
	return rs.redisClient.LLen(rs.ctx, message_cache_key).Val()
}
