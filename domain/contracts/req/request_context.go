package req

import (
	"context"

	"github.com/amahdian/ai-assistant-be/svc/auth"
	"github.com/gin-gonic/gin"
)

type RequestContext struct {
	Ctx      context.Context
	UserInfo *auth.UserInfo
}

func GetRequestContext(c *gin.Context) RequestContext {
	ctx := c.Request.Context()
	userInfo := auth.UserInfoFromCtx(ctx)

	return RequestContext{
		Ctx:      ctx,
		UserInfo: &userInfo,
	}
}
