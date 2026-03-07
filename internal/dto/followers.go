package dto

import "time"

type FollowersCreateDto struct {
	UserID     int64
	FollowerID int64
	CreatedAt  time.Time
}
