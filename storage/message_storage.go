package storage

import (
	"github.com/amahdian/ai-assistant-be/domain/model"
)

type MessageStorage interface {
	CrudStorage[*model.Message]

	ListByChatId(chatId string) ([]*model.Message, error)
}
