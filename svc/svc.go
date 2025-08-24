package svc

import (
	"context"
	"github.com/amahdian/ai-assistant-be/clients"

	"github.com/amahdian/ai-assistant-be/global/env"

	"github.com/amahdian/ai-assistant-be/storage"
)

type Svc interface {
	NewUserSvc(ctx context.Context) UserSvc
	NewChatSvc(ctx context.Context) ChatSvc
}

type svcImpl struct {
	stg       storage.Storage
	Envs      *env.Envs
	gptClient clients.GPTClient
}

func NewSvc(stg storage.Storage, envs *env.Envs, gptClient clients.GPTClient) Svc {
	return &svcImpl{
		stg,
		envs,
		gptClient,
	}
}

func (s *svcImpl) NewUserSvc(ctx context.Context) UserSvc {
	return newUserSvc(ctx, s.stg, s.Envs)
}

func (s *svcImpl) NewChatSvc(ctx context.Context) ChatSvc {
	return newChatSvc(ctx, s.stg, s.gptClient)
}
