package pg

import (
	"github.com/amahdian/ai-assistant-be/domain/model"
)

type ChatStg struct {
	crudStg[*model.Chat]
}

func NewChatStg(ses *ormSession) *ChatStg {
	return &ChatStg{
		crudStg: crudStg[*model.Chat]{db: ses.db},
	}
}

func (stg *ChatStg) ListByUserId(userId string) ([]*model.Chat, error) {
	var chats []*model.Chat
	err := stg.db.
		Where("user_id = ?", userId).
		Find(&chats).
		Error

	return chats, err
}
