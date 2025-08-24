package pg

import (
	"github.com/amahdian/ai-assistant-be/domain/model"
)

type MessageStg struct {
	crudStg[*model.Message]
}

func NewMessageStg(ses *ormSession) *MessageStg {
	return &MessageStg{
		crudStg: crudStg[*model.Message]{db: ses.db},
	}
}

func (stg *MessageStg) ListByChatId(chatId string) ([]*model.Message, error) {
	var messages []*model.Message
	err := stg.db.
		Where("chat_id = ?", chatId).
		Find(&messages).
		Error

	return messages, err
}
