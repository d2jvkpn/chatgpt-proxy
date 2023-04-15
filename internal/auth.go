package internal

import (
	// "fmt"
	"net/http"
	"strings"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/gin-gonic/gin"
)

func auth(ctx *gin.Context) {
	if !settings.AllowIps.Check(ctx.ClientIP()) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": -10, "msg": "unauthorized"})
		ctx.Abort()
		return
	}

	if !settings.AllowApiKeys.Enable {
		ctx.Next()
		return
	}

	// Authorization: Bearer
	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" || strings.HasPrefix(auth, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": -11, "msg": "unauthorized"})
		ctx.Abort()
	}

	if !settings.AllowApiKeys.Check(auth[7:]) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": -12, "msg": "unauthorized"})
		ctx.Abort()
	}

	ctx.Next()
}
