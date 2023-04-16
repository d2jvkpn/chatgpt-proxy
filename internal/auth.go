package internal

import (
	// "fmt"
	"net/http"
	"strings"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func auth(ctx *gin.Context) {
	ok := false

	defer func() {
		if ok {
			return
		}

		settings.ReqLogger.Error(
			"unauthorized access",
			zap.String("ip", ctx.ClientIP()),
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()),
		)
	}()

	if !settings.AllowIps.Check(ctx.ClientIP()) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": -10, "msg": "unauthorized"})
		ctx.Abort()
		return
	}

	if !settings.AllowApiKeys.Enable {
		ok = true
		ctx.Next()
		return
	}

	// Authorization: Bearer
	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" || strings.HasPrefix(auth, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": -11, "msg": "unauthorized"})
		ctx.Abort()
		return
	}

	if !settings.AllowApiKeys.Check(auth[7:]) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": -12, "msg": "unauthorized"})
		ctx.Abort()
		return
	}

	ok = true
	ctx.Next()
}
