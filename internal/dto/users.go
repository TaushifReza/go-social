package dto

import "time"

type UserResponseDto struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRegisterationDto struct {
	UserName string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
