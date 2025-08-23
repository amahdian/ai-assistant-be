package storage

import (
	"github.com/amahdian/ai-assistant-be/domain/model"
)

type UserStorage interface {
	CrudStorage[*model.User]

	FindByEmail(email string) (*model.User, error)
}
