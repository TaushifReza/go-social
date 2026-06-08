package cache

import (
	"context"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*dto.UserResponseDto, error)
		Set(context.Context, *dto.UserResponseDto) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}
