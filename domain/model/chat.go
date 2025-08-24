package model

import (
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId    string    `json:"user_id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`

	User     *User      `gorm:"-" json:"-"`
	Messages []*Message `gorm:"-" json:"messages,omitempty"`
}

func (*Chat) TableName() string {
	return "chats"
}
