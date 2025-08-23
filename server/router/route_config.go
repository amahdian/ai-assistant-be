package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type routeConfig struct {
	RequireUserSettings bool
	Middlewares         []gin.HandlerFunc
}

func newRouteConfig() *routeConfig {
	return &routeConfig{
		RequireUserSettings: false,
		Middlewares:         []gin.HandlerFunc{},
	}
}

func (rc *routeConfig) withUserSettings(flag bool) *routeConfig {
	clone := rc.clone()
	clone.RequireUserSettings = flag
	return clone

}

func (rc *routeConfig) withMiddlewares(middlewares ...gin.HandlerFunc) *routeConfig {
	clone := rc.clone()
	clone.Middlewares = append(rc.Middlewares, middlewares...)
	return clone
}

func (rc *routeConfig) withCompression() *routeConfig {
	clone := rc.clone()
	clone.Middlewares = append(rc.Middlewares, gzip.Gzip(gzip.DefaultCompression))
	return clone
}

func (rc *routeConfig) clone() *routeConfig {
	middlewares := make([]gin.HandlerFunc, len(rc.Middlewares))
	copy(middlewares, rc.Middlewares)

	return &routeConfig{
		Middlewares: middlewares,
	}
}
