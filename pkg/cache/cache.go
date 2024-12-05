package cache

import (
	"context"
	"errors"
	"fmt"
	redis2 "gin_template/internal/redis"
	"gin_template/pkg/logs"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var LocalCache *freecache.Cache

type StatusInfo struct {
	HitRate       string // 获取缓存命中率
	HitCount      int64  // 获取命中次数
	MissCount     int64  // 获取未命中次数
	EntryCount    int64  // 获取当前缓存条目数
	EvacuateCount int64  // 获取被清理的条目数
	TotalRequests int64  // 总请求次数
}

type KeyInfo struct {
	TTL   int64  // 过期时间
	Value string // 缓存的值
}

type Fetcher interface {
	Fetch(ctx context.Context, key string) (string, error)
}

type FetcherFunc func(ctx context.Context, key string) (string, error)

func (f FetcherFunc) Fetch(ctx context.Context, key string) (string, error) {
	return f(ctx, key)
}

func InitLocalCache() {
	// 创建一个 10MB 大小的缓存
	cacheSize := 10 * 1024 * 1024 // 10MB
	cache := freecache.NewCache(cacheSize)
	LocalCache = cache
}

// GetCacheStatus 封装的方法：获取缓存命中率、缓存命中数、总请求数
func GetCacheStatus() StatusInfo {
	hitRate := LocalCache.HitRate()             // 获取缓存命中率
	hitCount := LocalCache.HitCount()           // 获取命中次数
	missCount := LocalCache.MissCount()         // 获取未命中次数
	evacuateCount := LocalCache.EvacuateCount() // 获取被清理的条目数
	entryCount := LocalCache.EntryCount()       // 获取当前缓存条目数
	totalRequests := hitCount + missCount       // 计算总请求次数
	return StatusInfo{
		HitRate:       fmt.Sprintf("%.2f%%", hitRate*100),
		HitCount:      hitCount,
		MissCount:     missCount,
		TotalRequests: totalRequests,
		EntryCount:    entryCount,
		EvacuateCount: evacuateCount,
	}
}

// GetKeyStatus 封装的方法：获取对应缓存的值和过期时间
func GetKeyStatus(key string) KeyInfo {
	value, ttl, err := LocalCache.GetWithExpiration([]byte(key))
	if err != nil {
		logs.Log.Error("Error getting key:", zap.Error(err))
	}

	return KeyInfo{
		TTL:   int64(ttl) - time.Now().Unix(),
		Value: string(value),
	}
}

// GetLocal 获取本地缓存中的值
func GetLocal(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key is empty")
	}
	if LocalCache == nil {
		return "", nil
	}
	value, err := LocalCache.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// GetCache 先从本地缓存获取，若不存在，则从redis获取
func GetCache(ctx context.Context, key string) (string, error) {
	value, err := GetLocal(key)
	if err != nil {
		if !errors.Is(err, freecache.ErrNotFound) {
			return "", err
		}
	}
	if len(value) != 0 {
		return value, nil
	}
	pipeline := redis2.Client.Pipeline()
	getCmd := pipeline.Get(ctx, key)
	ttlCmd := pipeline.TTL(ctx, key)
	_, err = pipeline.Exec(ctx)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return "", err
		}
		return "", nil
	}
	value, err = getCmd.Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return "", err
		}
		return "", nil
	}
	ttl, err := ttlCmd.Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return "", err
		}
		return "", nil
	}
	duration := ttl / 3
	if int(duration.Seconds()) > 0 {
		err = SetLocal(key, value, duration)
		if err != nil {
			return "", errors.New("cache set err")
		}
	}
	return value, nil
}

// GetCacheOrElse 先从本地缓存获取，若不存在，则从redis获取，若redis也不存在，则调用fetcher
func GetCacheOrElse(ctx context.Context, key string, ttl time.Duration, fetcher Fetcher) (string, error) {
	value, err := GetLocal(key)
	if err != nil {
		if !errors.Is(err, freecache.ErrNotFound) {
			return "", err
		}
	}
	if len(value) != 0 {
		return value, nil
	}
	pipeline := redis2.Client.Pipeline()
	getCmd := pipeline.Get(ctx, key)
	ttlCmd := pipeline.TTL(ctx, key)
	_, err = pipeline.Exec(ctx)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return "", err
		}
	}
	remainingTime, err := ttlCmd.Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return "", err
		}
	}
	value, err = getCmd.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			value, err = fetcher.Fetch(ctx, key)
			if err != nil {
				return "", err
			}
			SetCache(ctx, key, value, ttl)
			return value, nil
		} else {
			return "", err
		}
	}
	duration := remainingTime / 3
	if int(duration.Seconds()) > 0 {
		err = SetLocal(key, value, duration)
		if err != nil {
			return "", errors.New("cache set err")
		}
	}
	return value, nil
}

// SetLocal 设置本地缓存值
func SetLocal(key, value string, ttl time.Duration) error {
	if len(key) == 0 {
		return errors.New("key is empty")
	}
	if LocalCache == nil {
		return nil
	}
	err := LocalCache.Set([]byte(key), []byte(value), int(ttl.Seconds()))
	if err != nil {
		logs.Log.Error("cache set err:", zap.Error(err))
		return errors.New("cache set err")
	}
	return nil
}

// SetCache 设置redis和本地缓存值
func SetCache(ctx context.Context, key, value string, ttl time.Duration) error {
	if len(key) == 0 {
		return errors.New("key is empty")
	}
	duration := ttl / 3
	if int(duration.Seconds()) > 0 {
		err := SetLocal(key, value, duration)
		if err != nil {
			return err
		}
	}
	result, err := redis2.Client.Set(ctx, key, value, ttl).Result()
	logs.Log.Info("redis2 set", zap.String("result", result))
	if err != nil {
		return err
	}
	return nil
}

// DelLocal 删除本地缓存值
func DelLocal(key string) error {
	if len(key) == 0 {
		return errors.New("key is empty")
	}
	if LocalCache == nil {
		return nil
	}
	LocalCache.Del([]byte(key))
	return nil
}

// DelCache 删除redis和本地缓存值
func DelCache(key string) error {
	if len(key) == 0 {
		return errors.New("key is empty")
	}
	err := DelLocal(key)
	if err != nil {
		return err
	}
	result, err := redis2.Client.Del(context.Background(), key).Result()
	logs.Log.Info("redis2 del num", zap.Int64("result", result))
	if err != nil {
		return err
	}
	return nil
}
