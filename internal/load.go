package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/d2jvkpn/chatgpt-proxy/internal/handlers"
	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/d2jvkpn/x-ai/pkg/chatgpt"
	"github.com/gin-gonic/gin"
)

func Load(config string, release bool) (err error) {
	var (
		engine *gin.Engine
		router *gin.RouterGroup
	)

	//
	if settings.GPTCli, err = chatgpt.NewClient(config, "chatgpt"); err != nil {
		return err
	}

	if settings.AllowIps, err = settings.NewAllowedKeys(config, "allow_ips"); err != nil {
		return err
	}

	if settings.AllowApiKeys, err = settings.NewAllowedKeys(config, "allow_api_keys"); err != nil {
		return err
	}

	if settings.Tls, err = settings.NewTlsConfig(config, "tls"); err != nil {
		return err
	}

	if settings.AllowApiKeys.Enable && !settings.Tls.Enable {
		return fmt.Errorf("enabled api keys without using tls")
	}

	//
	if release {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(gin.Recovery())
	} else {
		engine = gin.Default()
	}
	engine.Use(cors("*"))

	router = &engine.RouterGroup
	handlers.RouteOpen(router)
	handlers.RouteChatgpt(router, auth)

	_Server.Handler = engine

	return nil
}

func cors(origin string) gin.HandlerFunc {
	allowHeaders := strings.Join([]string{"Content-Type", "Authorization"}, ", ")

	exposeHeaders := strings.Join([]string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Content-Type",
		"Content-Length",
	}, ", ")

	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", origin)

		ctx.Header("Access-Control-Allow-Headers", allowHeaders)
		// Content-Type, Authorization, X-CSRF-Token
		ctx.Header("Access-Control-Expose-Headers", exposeHeaders)
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
}
