package model

import (
	"github.com/amahdian/ai-assistant-be/domain/model/common"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	ChatID    string          `json:"chat_id"`
	Role      string          `json:"role"` // "user" or "assistant"
	Content   string          `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	Metadata  common.Metadata `json:"metadata" gorm:"type:jsonb"`

	Chat *Chat `gorm:"-" json:"chat,omitempty"`
}

func (*Message) TableName() string {
	return "messages"
}

type StreamedMessage struct {
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
