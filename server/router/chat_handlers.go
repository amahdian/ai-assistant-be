package router

import (
	"github.com/amahdian/ai-assistant-be/domain/contracts/req"
	"github.com/amahdian/ai-assistant-be/domain/contracts/resp"
	"github.com/gin-gonic/gin"
	"io"
)

func (r *Router) listChats(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	chatSvc := r.svc.NewChatSvc(ctx)

	res, err := chatSvc.ListChats(&user)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}
	resp.Ok(ctx, res)
}

func (r *Router) getChat(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	var reqUri req.IdUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewChatSvc(reqCtx.Ctx)
	res, err := dSvc.GetChat(reqUri.Id, &user)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, res)
}

func (r *Router) deleteChat(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	var reqUri req.IdUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewChatSvc(reqCtx.Ctx)
	err := dSvc.DeleteChat(reqUri.Id, &user)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, true)
}

func (r *Router) createChat(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	request := &req.SendMessage{}
	err := ctx.BindJSON(&request)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewChatSvc(reqCtx.Ctx)
	res, err := dSvc.CreateChat(request.Message, &user)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, res)
}

func (r *Router) sendMessage(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)
	user := reqCtx.UserInfo.User()

	var reqUri req.IdUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	request := &req.SendMessage{}
	if err := ctx.BindJSON(&request); err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	chatSvc := r.svc.NewChatSvc(reqCtx.Ctx)

	// Check for the 'stream' query parameter
	useStream := ctx.DefaultQuery("stream", "false") == "true"

	if useStream {
		// --- Streaming Response ---
		streamChan, err := chatSvc.SendMessageStream(reqUri.Id, request.Message, &user)
		if err != nil {
			resp.AbortWithError(ctx, err)
			return
		}

		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		ctx.Header("Access-Control-Allow-Origin", "*")

		ctx.Stream(func(w io.Writer) bool {
			if msg, ok := <-streamChan; ok {
				ctx.SSEvent("message", gin.H{
					"content":  msg.Content,
					"metadata": msg.Metadata,
				})
				return true
			}
			return false
		})
	} else {
		// --- Non-streaming (standard JSON) Response ---
		res, err := chatSvc.SendMessage(reqUri.Id, request.Message, &user)
		if err != nil {
			resp.AbortWithError(ctx, err)
			return
		}
		resp.Ok(ctx, res)
	}
}
