package storage

import (
	"github.com/amahdian/ai-assistant-be/domain/model"
)

type ChatStorage interface {
	CrudStorage[*model.Chat]

	ListByUserId(userId string) ([]*model.Chat, error)
}
