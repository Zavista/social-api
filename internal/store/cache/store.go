package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/zavista/social-api/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.UserWithRole, error)
		Set(context.Context, *store.UserWithRole) error
		Delete(context.Context, int64)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
	}
}
