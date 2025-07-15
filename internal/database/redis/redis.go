package redis

import (
	"CloudPhoto/config"
	"CloudPhoto/internal/tool"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Get().Redis.Host + ":" + config.Get().Redis.Port,
		Password: config.Get().Redis.Password,
		DB:       config.Get().Redis.DB,
	})
}

func Get(ctx *context.Context, key string) string {
	val, err := rdb.Get(*ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}
func Set(ctx *context.Context, key string, value any, ttl time.Duration) {
	tool.PanicIfErr(rdb.Set(*ctx, key, value, ttl).Err())
}

func Del(ctx *context.Context, key string) {
	tool.PanicIfErr(rdb.Del(*ctx, key).Err())
}
