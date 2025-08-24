package router

import (
	"net/http"

	"github.com/amahdian/ai-assistant-be/server/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupRoutes() {
	r.publicGroup = r.Group("")
	r.authGroup = r.Group(
		"",
		middleware.VerifyAuth(r.authenticator),
	)

	r.registerPublicRoutes()
	r.registerUserRoutes()
	r.registerChatRoutes()
}

func (r *Router) registerPublicRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.publicGroup, http.MethodGet, "/health", r.healthCheck, config)
	r.registerRoute(r.publicGroup, http.MethodGet, "/swagger/*any", r.swaggerHandler, config)
}

func (r *Router) registerUserRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.publicGroup, http.MethodPost, "/user/login", r.login, config)
	r.registerRoute(r.publicGroup, http.MethodPost, "/user/register", r.register, config)
}

func (r *Router) registerChatRoutes() {
	config := newRouteConfig()
	r.registerRoute(r.authGroup, http.MethodGet, "/chat", r.listChats, config)
	r.registerRoute(r.authGroup, http.MethodGet, "/chat/:id", r.getChat, config)
	r.registerRoute(r.authGroup, http.MethodDelete, "/chat/:id", r.deleteChat, config)
	r.registerRoute(r.authGroup, http.MethodPost, "/chat", r.createChat, config)
	r.registerRoute(r.authGroup, http.MethodPost, "/chat/:id", r.sendMessage, config)
}

func (r *Router) registerRoute(routerGroup *gin.RouterGroup, method, path string, handler gin.HandlerFunc, configs ...*routeConfig) {
	config := newRouteConfig()
	if len(configs) > 0 {
		config = configs[0]
	}

	handlers := make([]gin.HandlerFunc, 0)

	if r.storage != nil && config.RequireUserSettings {
		handlers = append(handlers, middleware.WithUserSettings(r.storage))
	}

	if len(config.Middlewares) > 0 {
		handlers = append(handlers, config.Middlewares...)
	}

	handlers = append(handlers, handler)
	routerGroup.Handle(method, path, handlers...)
}
