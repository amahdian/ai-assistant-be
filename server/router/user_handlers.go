package router

import (
	"github.com/amahdian/ai-assistant-be/domain/contracts/req"
	"github.com/amahdian/ai-assistant-be/domain/contracts/resp"
	"github.com/gin-gonic/gin"
)

// login user login via jwt.
//
//	@Summary	login and get the jwt auth token
//	@Description
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		request	body		req.Login	true	"login credentials"
//	@Success	200		{object}	resp.Response[string]
//	@Failure	400		{object}	resp.ErrorResponse
//	@Failure	500		{object}	resp.ErrorResponse
//	@Security	Bearer
//	@Router		/api/v1/user/login [post]
func (r *Router) login(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	request := &req.Login{}
	err := ctx.BindJSON(&request)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewUserSvc(reqCtx.Ctx)
	res, err := dSvc.Login(request.Email, request.Password)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, res)
}

// register create a new user.
//
//	@Summary	register new user
//	@Description
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		request	body		req.Register	true	"register data"
//	@Success	200		{object}	resp.Response[string]
//	@Failure	400		{object}	resp.ErrorResponse
//	@Failure	500		{object}	resp.ErrorResponse
//	@Security	Bearer
//	@Router		/api/v1/user/register [post]
func (r *Router) register(ctx *gin.Context) {
	reqCtx := req.GetRequestContext(ctx)

	request := &req.Register{}
	err := ctx.BindJSON(&request)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	dSvc := r.svc.NewUserSvc(reqCtx.Ctx)
	res, err := dSvc.Register(request.Email, request.Password)
	if err != nil {
		resp.AbortWithError(ctx, err)
		return
	}

	resp.Ok(ctx, res)
}
