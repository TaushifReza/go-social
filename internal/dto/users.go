package dto

import (
	"time"

	"github.com/TaushifReza/go-social/internal/model"
)

type UserResponseDto struct {
	ID        int64      `json:"id"`
	UserName  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	RoleID    int64      `json:"role_id"`
	Role      model.Role `json:"role"`
}

type UserRegisterationDto struct {
	UserName string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleID   int64  `json:"role_id" validate:"required"`
}

type LoginDto struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRepoResponseDto struct {
	ID       int64
	Username string
	Email    string
	Password string
}
