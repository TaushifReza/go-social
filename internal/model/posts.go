package model

import "time"

type Posts struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
	Comment   []Comment `json:"comments"`
	User      User      `json:"user"`
}
