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
	if rdb == nil {
		return Storage{
			Users: &NoOpStore{}, // Safely returns a no-op implementation
		}
	}

	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}

type NoOpStore struct{}

func (n *NoOpStore) Get(ctx context.Context, userID int64) (*dto.UserResponseDto, error) {
	return nil, nil // Always act like a cache miss safely
}

func (n *NoOpStore) Set(ctx context.Context, user *dto.UserResponseDto) error {
	return nil // Safely do nothing
}
