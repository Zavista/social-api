package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/zavista/social-api/internal/store"
)

type CacheStorage struct {
	Users interface {
		Get(context.Context, int64) (*store.UserWithRole, error)
		Set(context.Context, *store.UserWithRole) error
		Delete(context.Context, int64)
	}
}

func NewRedisStorage(rbd *redis.Client) CacheStorage {
	return CacheStorage{
		Users: &UserStore{rdb: rbd},
	}
}
