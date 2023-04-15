package handlers

import (
	"net/http"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/gin-gonic/gin"
)

func RouteOpen(router *gin.RouterGroup, handers ...gin.HandlerFunc) {
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	open := router.Group("/api/v1/open", handers...)

	open.GET("/version", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok", "data": gin.H{"version": settings.Version()},
		})
	})
}

func RouteChatgpt(router *gin.RouterGroup, handers ...gin.HandlerFunc) {
	group := router.Group("/v1", handers...)

	group.POST("/chat/completions", completions)
	group.POST("/images/generations", imggen)
	group.POST("/models", models)
}
