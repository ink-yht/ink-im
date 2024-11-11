package user_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"ink-im-server/internal/domain/user_domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id uint) (user_domain.User, error)
	Set(ctx context.Context, u user_domain.User) error
}

type RedisUserCache struct {

	// 面向接口编程
	client redis.Cmdable
	// 超时时间
	expiration time.Duration
}

// A 用到了 B，B 一定是接口
// A 用到了 B，B 一定是 A 的字段
// A 用到了 B，A 绝不初始化 B，而是外面注入

// NewUserCache 依赖注入
func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 只要 error 为 nil，就认为缓存里一定有数据
// 如果没有数据，返回一个特定的 err

func (cache *RedisUserCache) Get(ctx context.Context, id uint) (user_domain.User, error) {
	key := cache.Key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return user_domain.User{}, err
	}
	var u user_domain.User
	err = json.Unmarshal(val, &u)
	if err != nil {
		return user_domain.User{}, err
	}
	return u, nil
}

func (cache *RedisUserCache) Set(ctx context.Context, u user_domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}

	key := cache.Key(u.Id)

	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) Key(id uint) string {
	return fmt.Sprintf("user:info:%d", id)
}
