package handlers

import (
	// "net/http"
	"log"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/gin-gonic/gin"
)

func langChainIndex(ctx *gin.Context) {
	var err error

	if _, err = settings.LCA.HandleIndex(ctx); err != nil {
		log.Printf("!!! langChainIndex: %v\n", err)
	}
}

func langChainQuery(ctx *gin.Context) {
	var err error

	if err = settings.LCA.HandleQuery(ctx); err != nil {
		log.Printf("!!! langChainQuery: %v\n", err)
	}
}
