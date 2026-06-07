package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	RoleID    int64     `json:"role_id"`
}

type UserInvitations struct {
	Token  string
	UserID int64
	Expire time.Time
}
