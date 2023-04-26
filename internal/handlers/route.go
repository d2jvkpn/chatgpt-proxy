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
			"code": 0, "msg": "ok", "data": gin.H{"version": settings.GetVersion()},
		})
	})
}

func RouteChatgpt(router *gin.RouterGroup, handers ...gin.HandlerFunc) {
	group := router.Group("/v1", handers...)
	group.GET("/models", listModels)
	group.POST("/completions", createCompletions)

	group.POST("/chat/completions", createChatCompletions)

	group.POST("/images/generations", createImage)
	group.POST("/images/edits", createEditImage)
	group.POST("/images/variations", createVariImage)

	group.POST("/embeddings", createEmbeddings)
}

func RouteAuth(router *gin.RouterGroup, handers ...gin.HandlerFunc) {
	group := router.Group("/api/v1/auth", handers...)

	langchain := group.Group("/langchain")
	langchain.POST("/index", langChainIndex)
	langchain.POST("/query", langChainQuery)
}
