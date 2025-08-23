package svc

import (
	"context"
	"errors"
	"time"

	"github.com/amahdian/ai-assistant-be/domain/model"
	"github.com/amahdian/ai-assistant-be/global/env"
	"github.com/amahdian/ai-assistant-be/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserSvc interface {
	Login(email, password string) (string, error)
	Register(email, password string) (string, error)
}

type userSvc struct {
	ctx context.Context
	stg storage.Storage

	envs *env.Envs
}

func newUserSvc(ctx context.Context, stg storage.Storage, envs *env.Envs) UserSvc {
	return &userSvc{
		ctx:  ctx,
		stg:  stg,
		envs: envs,
	}
}

func (s *userSvc) Login(email, password string) (string, error) {
	user, err := s.stg.User(s.ctx).FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(s.envs.Server.JwtSecret))

	return tokenStr, nil
}

func (s *userSvc) Register(email, password string) (string, error) {
	user, err := s.stg.User(s.ctx).FindByEmail(email)
	if err != nil {
		return "", err
	}

	if user != nil {
		return "", errors.New("user is already registered")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user = &model.User{
		Email:        email,
		PasswordHash: string(hash),
	}
	err = s.stg.User(s.ctx).CreateOne(user)
	if err != nil {
		return "", err
	}

	return s.Login(email, password)
}
