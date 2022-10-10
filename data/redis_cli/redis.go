package redis_cli

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisDB *redis.ClusterClient
)

func Init(cfg *Cfg) {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.Host,
		Password: cfg.Password,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.PoolTimeout))
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("连接Redis数据库失败，%s", err)
	}
	RedisDB = client
}
