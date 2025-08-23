package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/amahdian/ai-assistant-be/domain/model"
	"github.com/amahdian/ai-assistant-be/global/env"
	"github.com/amahdian/ai-assistant-be/global/errs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type userInfoCtx struct{}

type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func (u *UserInfo) User() model.User {
	return model.User{
		ID:    u.ID,
		Email: u.Email,
	}
}

type Authenticator interface {
	Verify(request *http.Request) (context.Context, error)
}

type authenticator struct {
	JwtSecret string
}

func NewAuthenticator(envs *env.Envs) Authenticator {
	return &authenticator{
		JwtSecret: envs.Server.JwtSecret,
	}
}

func (a *authenticator) Verify(request *http.Request) (context.Context, error) {
	ctx := request.Context()
	tokenStr := request.Header.Get("Authorization")
	if tokenStr == "" {
		return ctx, errors.New("authorization header is empty")
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.JwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx = context.WithValue(ctx, userInfoCtx{}, UserInfo{
			ID:    uuid.MustParse(claims["id"].(string)),
			Email: claims["email"].(string),
		})
		return ctx, nil
	} else {
		return ctx, errs.Newf(errs.Unauthenticated, err, "Auth failed.")
	}
}

func UserInfoFromCtx(ctx context.Context) UserInfo {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}
	u, ok := ctx.Value(userInfoCtx{}).(UserInfo)
	if !ok {
		return UserInfo{}
	}
	return u
}
