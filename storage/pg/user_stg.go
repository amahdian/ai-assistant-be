package pg

import (
	"errors"

	"github.com/amahdian/ai-assistant-be/domain/model"
	"gorm.io/gorm"
)

type UserStg struct {
	crudStg[*model.User]
}

func NewUserStg(ses *ormSession) *UserStg {
	return &UserStg{
		crudStg: crudStg[*model.User]{db: ses.db},
	}
}

func (stg *UserStg) FindByEmail(email string) (user *model.User, err error) {
	err = stg.db.
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return
}
