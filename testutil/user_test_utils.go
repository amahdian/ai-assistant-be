package testutil

import (
	"github.com/amahdian/ai-assistant-be/domain/model"
	"github.com/amahdian/ai-assistant-be/global/test"
)

func TestUser() *model.User {
	return &model.User{
		Email: test.UserEmail,
	}
}
