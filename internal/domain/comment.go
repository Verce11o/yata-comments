package domain

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	CommentID uuid.UUID `json:"comment_id" db:"comment_id"`
	TweetID   uuid.UUID `json:"tweet_id" db:"tweet_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Text      string    `json:"text" db:"text"`
	ImageName string    `json:"image_name"`
	ImageURL  string    `json:"image_url" db:"image_temp_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
