package dto

import (
	"time"
)

type CreatePostDto struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostDto struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

type CommentDto struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type UserDto struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
}

type PostsDto struct {
	ID        int64        `json:"id"`
	Content   string       `json:"content"`
	Title     string       `json:"title"`
	UserID    int64        `json:"user_id"`
	Tags      []string     `json:"tags"`
	CreatedAt time.Time    `json:"created_at"`
	Version   int64        `json:"version"`
	UpdatedAt time.Time    `json:"updated_at"`
	Comment   []CommentDto `json:"comments"`
	User      UserDto      `json:"user"`
}

type PostWithMetaData struct {
	PostsDto
	CommentCount int `json:"comments_count"`
}
