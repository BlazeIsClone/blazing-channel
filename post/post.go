package post

import (
	"errors"
	"time"
)

type Post struct {
	ID        int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePostCommand struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpdatePostCommand struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

var ErrNotFound = errors.New("post not found")
