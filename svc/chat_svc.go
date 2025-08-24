package svc

import (
	"context"
	"errors"
	"fmt"
	"github.com/amahdian/ai-assistant-be/clients"
	"github.com/amahdian/ai-assistant-be/domain/model"
	"github.com/amahdian/ai-assistant-be/domain/model/common"
	"github.com/amahdian/ai-assistant-be/global/errs"
	"github.com/amahdian/ai-assistant-be/pkg/logger"
	"github.com/amahdian/ai-assistant-be/storage"
)

// ChatSvc defines the interface for chat-related services.
type ChatSvc interface {
	DeleteChat(chatID string, user *model.User) error
	CreateChat(message string, user *model.User) (*model.Chat, error)
	SendMessage(chatID, message string, user *model.User) (*model.Message, error)
	SendMessageStream(chatID, message string, user *model.User) (<-chan *model.StreamedMessage, error)
	ListChats(user *model.User) ([]*model.Chat, error)
	GetChat(id string, user *model.User) (*model.Chat, error)
}

type chatSvc struct {
	ctx       context.Context
	stg       storage.Storage
	gptClient clients.GPTClient
}

func newChatSvc(ctx context.Context, stg storage.Storage, gptClient clients.GPTClient) ChatSvc {
	return &chatSvc{
		ctx:       ctx,
		stg:       stg,
		gptClient: gptClient,
	}
}

func (s *chatSvc) SendMessage(chatID, message string, user *model.User) (*model.Message, error) {
	// 1. Validate chat ownership and save user message
	chat, err := s.prepareMessage(chatID, message, user)
	if err != nil {
		return nil, err
	}

	// 2. Get conversation history
	messages, err := s.stg.Message(s.ctx).ListByChatId(chatID)
	if err != nil {
		return nil, errs.Wrapf(err, "failed to list messages")
	}

	agent := model.DefaultAgent

	// 3. Check for summarization
	chatSummary := s.checkAndSummarizeIfNeeded(messages)

	var reply string

	reply, err = s.gptClient.SendToGPT(agent.SystemPrompt, chatSummary, messages)
	if err != nil {
		return nil, errs.Wrapf(err, "failed to get GPT response")
	}

	// 5. Save the assistant's message
	assistantMessage := &model.Message{
		ChatID:  chatID,
		Role:    "assistant",
		Content: reply,
		Chat:    chat,
	}

	if err = s.stg.Message(s.ctx).CreateOne(assistantMessage); err != nil {
		return nil, errs.Wrapf(err, "failed to save assistant message")
	}

	return assistantMessage, nil
}

func (s *chatSvc) SendMessageStream(chatID, message string, user *model.User) (<-chan *model.StreamedMessage, error) {
	_, err := s.prepareMessage(chatID, message, user)
	if err != nil {
		return nil, err
	}

	messages, err := s.stg.Message(s.ctx).ListByChatId(chatID)
	if err != nil {
		return nil, errs.Wrapf(err, "failed to list messages")
	}

	agent := model.DefaultAgent
	chatSummary := s.checkAndSummarizeIfNeeded(messages)

	// Streaming mode for other agents
	stream, err := s.gptClient.SendToGPTStream(agent.SystemPrompt, chatSummary, messages)
	if err != nil {
		return nil, errs.Wrapf(err, "failed to start GPT stream")
	}

	msgChan := make(chan *model.StreamedMessage)
	go func() {
		defer close(msgChan)
		var fullReply string
		for chunk := range stream {
			fullReply += chunk
			msgChan <- &model.StreamedMessage{Content: chunk}
		}

		if fullReply != "" {
			_ = s.stg.Message(s.ctx).CreateOne(&model.Message{
				ChatID:  chatID,
				Role:    "assistant",
				Content: fullReply,
			})
		}
	}()

	return msgChan, nil
}

func (s *chatSvc) ListChats(user *model.User) ([]*model.Chat, error) {
	return s.stg.Chat(s.ctx).ListByUserId(user.ID.String())
}

func (s *chatSvc) GetChat(id string, user *model.User) (*model.Chat, error) {
	chat, err := s.stg.Chat(s.ctx).FindById(id)
	if err != nil {
		return nil, err
	}
	if chat.UserId != user.ID.String() {
		return nil, errors.New("permission denied")
	}
	messages, err := s.stg.Message(s.ctx).ListByChatId(id)
	if err != nil {
		return nil, err
	}
	chat.Messages = messages
	return chat, nil
}

func (s *chatSvc) DeleteChat(chatID string, user *model.User) error {
	chat, err := s.stg.Chat(s.ctx).FindById(chatID)
	if err != nil {
		return err
	}
	if chat.UserId != user.ID.String() {
		return errors.New("not authorized")
	}
	return s.stg.Chat(s.ctx).DeleteById(chatID)
}

func (s *chatSvc) CreateChat(message string, user *model.User) (*model.Chat, error) {
	title, err := s.createChatTitle(message)
	if err != nil {
		return nil, errs.Wrapf(err, "failed to create chat title")
	}
	newChat := model.Chat{
		Title:  title,
		UserId: user.ID.String(),
	}
	if err = s.stg.Chat(s.ctx).CreateOne(&newChat); err != nil {
		return nil, errs.Wrapf(err, "failed to create chat record")
	}
	return &newChat, nil
}

func (s *chatSvc) prepareMessage(chatID, message string, user *model.User) (*model.Chat, error) {
	chat, err := s.stg.Chat(s.ctx).FindById(chatID)
	if err != nil {
		return nil, errs.Wrapf(err, "failed to find chat")
	}
	if chat.UserId != user.ID.String() {
		return nil, errors.New("permission denied")
	}

	err = s.stg.Message(s.ctx).CreateOne(&model.Message{
		ChatID:   chatID,
		Role:     "user",
		Content:  message,
		Metadata: common.Metadata{},
	})
	if err != nil {
		return nil, errs.Wrapf(err, "failed to save user message")
	}
	return chat, nil
}

func (s *chatSvc) createChatTitle(message string) (string, error) {
	titlePrompt := []*model.Message{
		{Role: "system", Content: "You are a helpful assistant that writes concise titles for chat conversations."},
		{Role: "user", Content: fmt.Sprintf("Create a short and clear title without quoutes for this message:\n\"%s\"", message)},
	}
	return s.gptClient.SendMessages(titlePrompt)
}

func (s *chatSvc) checkAndSummarizeIfNeeded(messages []*model.Message) string {
	if len(messages) <= 40 {
		return ""
	}
	var chunks []string
	for _, m := range messages {
		chunks = append(chunks, fmt.Sprintf("%s: %s", m.Role, m.Content))
	}
	textToSummarize := "Summarize this conversation in 3 lines:\n" + joinChunks(chunks, 3000)
	summaryMessages := []*model.Message{
		{Role: "system", Content: "You are a helpful summarizer."},
		{Role: "user", Content: textToSummarize},
	}
	summary, err := s.gptClient.SendMessages(summaryMessages)
	if err != nil {
		logger.Debugf("could not summarize chat: %v\n", err)
		return ""
	}
	return summary
}

func joinChunks(chunks []string, limit int) string {
	var result string
	count := 0
	for _, s := range chunks {
		if count+len(s) > limit {
			break
		}
		result += s + "\n"
		count += len(s) + 1
	}
	return result
}
