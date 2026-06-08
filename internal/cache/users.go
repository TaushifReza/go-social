package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/redis/go-redis/v9"
)

const UserExpTime = time.Minute

type UserStore struct {
	rdb *redis.Client
}

func (s *UserStore) Get(ctx context.Context, userID int64) (*dto.UserResponseDto, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var user dto.UserResponseDto
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *dto.UserResponseDto) error {
	// Guard clause against nil pointer panics
	if user == nil {
		return fmt.Errorf("cannot cache nil user")
	}
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)

	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, json, UserExpTime).Err()

}
