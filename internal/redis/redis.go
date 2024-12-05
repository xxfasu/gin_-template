package redis

import (
	"context"
	"gin_template/internal/conf"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"

	"time"
)

var Client *redis.Client
var RLock *redsync.Redsync

func InitRedis() error {
	// 创建 Redis 客户端
	Client = redis.NewClient(&redis.Options{
		Addr:         conf.Config.Redis.Addr,
		Password:     conf.Config.Redis.Password, // 如果没有密码，使用空字符串
		DB:           conf.Config.Redis.DB,
		Username:     conf.Config.Redis.Username,
		PoolSize:     16, // 连接池大小
		MinIdleConns: 3,  // 最小空闲连接数
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()

	if err != nil {
		return err
	}
	// 创建 Redsync 的 Redis 连接池
	pool := goredis.NewPool(Client)

	// 创建 Redsync 实例
	RLock = redsync.New(pool)
	return nil
}
