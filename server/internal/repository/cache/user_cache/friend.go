package user_cache

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type FriendCache interface {
	//Get(ctx context.Context, id uint) (user_domain.User, error)
	//Set(ctx context.Context, u user_domain.User) error
}

type RedisFriendCache struct {

	// 面向接口编程
	client redis.Cmdable
	// 超时时间
	expiration time.Duration
}

// A 用到了 B，B 一定是接口
// A 用到了 B，B 一定是 A 的字段
// A 用到了 B，A 绝不初始化 B，而是外面注入

// NewFriendCache 依赖注入
func NewFriendCache(client redis.Cmdable) FriendCache {
	return &RedisFriendCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}
