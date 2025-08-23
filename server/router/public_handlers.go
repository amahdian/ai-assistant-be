package router

import (
	"net/http"
	"time"

	"github.com/amahdian/ai-assistant-be/domain/contracts/resp"
	"github.com/amahdian/ai-assistant-be/version"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ginSwaggerHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)

// healthCheck returns health status of the server
//
//	@Summary	health check
//	@Description
//	@Tags		Public
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	resp.HealthResponseDto
//	@Router		/health [get]
func (r *Router) healthCheck(ctx *gin.Context) {
	response := resp.HealthResponseDto{
		AppName:    version.AppName,
		AppVersion: version.AppVersion,
	}

	resp.Ok(ctx, response)
}

func (r *Router) swaggerHandler(ctx *gin.Context) {
	// the ginSwaggerHandler by default recognizes the "/swagger/index.html"  but not"/swagger" or "/swagger/".
	// therefore we add support for these endpoints by redirecting to "/swagger/index.html"
	if ctx.Request.RequestURI == "/swagger" || ctx.Request.RequestURI == "/swagger/" {
		ctx.Redirect(http.StatusFound, "/swagger/index.html")
	}
	ginSwaggerHandler(ctx)
}

// getServerTime returns time of the server
//
//	@Summary	getServerTime
//	@Description
//	@Tags		Public
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	string
//	@Router		/server-time [get]
func (r *Router) getServerTime(ctx *gin.Context) {
	resp.Ok(ctx, time.Now())
}
