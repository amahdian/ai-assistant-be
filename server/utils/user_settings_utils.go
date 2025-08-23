package utils

import (
	"context"

	"github.com/amahdian/ai-assistant-be/domain/model"
	"github.com/gin-gonic/gin"
)

func CurrentUserSettings(ctx context.Context) *model.User {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}
	us, ok := ctx.Value(UserSettingsContextKey).(*model.User)
	if !ok {
		return nil
	}
	return us
}
