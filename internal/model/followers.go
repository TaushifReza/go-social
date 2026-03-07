package model

import "time"

type Followers struct {
	UserID     int64
	FollowerID int64
	CreatedAt  time.Time
}
