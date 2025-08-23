package resp

import (
	"github.com/amahdian/ai-assistant-be/domain/model/common"
	"github.com/amahdian/ai-assistant-be/pkg/msg"
)

type ErrorResponse struct {
	Success  bool                      `json:"success" default:"false"`
	Error    string                    `json:"error,omitempty"`
	Messages map[string][]*msg.Message `json:"messages,omitempty"`
}

type Response[T any] struct {
	Success  bool                      `json:"success" default:"true"`
	Data     T                         `json:"data,omitempty"`
	Messages map[string][]*msg.Message `json:"messages,omitempty"`
}

type PaginatedResponse[T any] struct {
	Success  bool     `json:"success"`
	Data     []T      `json:"data"`
	PageInfo PageInfo `json:"pageInfo"`
}

type PageInfo struct {
	Page          int   `json:"page"`
	PageSize      int   `json:"pageSize"`
	ElementsCount int   `json:"elementsCount"`
	TotalCount    int64 `json:"totalCount"`
	HasMore       bool  `json:"hasMore"`
	IsEmpty       bool  `json:"isEmpty"`
}

type HealthResponseDto struct {
	AppName    string `json:"appName"`
	AppVersion string `json:"appVersion"`
}

func NewErrorResponse(data any) *ErrorResponse {
	resp := &ErrorResponse{
		Success: false,
	}
	switch v := data.(type) {
	case error:
		resp.Error = msg.MakePlainText(v.Error())
	case string:
		resp.Error = msg.MakePlainText(v)
	case *msg.MessageContainer:
		resp.Messages = v.GetAll()
	}
	return resp
}

func NewResponse[T any](data T) *Response[T] {
	return &Response[T]{
		Success: true,
		Data:    data,
	}
}

func NewResponseWithMessage[T any](data T, messages *msg.MessageContainer) *Response[T] {
	return &Response[T]{
		Success:  true,
		Data:     data,
		Messages: messages.GetAll(),
	}
}

func NewPaginatedResponse[T any](data []T, pagination *common.Pagination) *PaginatedResponse[T] {
	if data == nil {
		data = make([]T, 0)
	}
	elementsCount := len(data)
	hasMore := int64(pagination.Page)*int64(pagination.PageSize)+int64(elementsCount) < pagination.TotalCount
	isEmpty := elementsCount == 0
	return &PaginatedResponse[T]{
		Success: true,
		Data:    data,
		PageInfo: PageInfo{
			Page:          pagination.Page,
			PageSize:      pagination.PageSize,
			ElementsCount: elementsCount,
			TotalCount:    pagination.TotalCount,
			HasMore:       hasMore,
			IsEmpty:       isEmpty,
		},
	}
}
