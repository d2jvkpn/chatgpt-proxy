package handlers

import (
	// "net/http"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/gin-gonic/gin"
)

func langChainIndex(ctx *gin.Context) {
	_, _ = settings.LCA.HandleIndex(ctx)
}

func langChainQuery(ctx *gin.Context) {
	_ = settings.LCA.HandleQuery(ctx)
}
