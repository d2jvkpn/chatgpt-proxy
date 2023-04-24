package internal

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/d2jvkpn/chatgpt-proxy/internal/biz"
	"github.com/d2jvkpn/chatgpt-proxy/internal/handlers"
	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/d2jvkpn/go-web/pkg/wrap"
	xwrap "github.com/d2jvkpn/x-ai/pkg/wrap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Load(config string, release bool) (err error) {
	var (
		vp     *viper.Viper
		engine *gin.Engine
		router *gin.RouterGroup
	)

	vp = viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigFile(config)
	if err = vp.ReadInConfig(); err != nil {
		return err
	}

	//
	/*
		if settings.GPTCli, err = chatgpt.NewClient(config, "chatgpt"); err != nil {
			return err
		}
	*/

	if settings.GPTCli2, err = xwrap.NewOpenAiClient(config, "chatgpt"); err != nil {
		return err
	}

	settings.LCA, err = biz.NewLangChainAgent(vp.GetString("chatgpt.api_key"), "data/langchain")
	if err != nil {
		return
	}

	//
	if _Tls, err = NewTlsConfig(config, "tls"); err != nil {
		return err
	}

	//
	level := zap.DebugLevel
	if release {
		level = zap.InfoLevel
	}
	settings.Logger = wrap.NewLogger("logs/chatgot-proxy.log", level, 256, nil)
	settings.SetupLoggers()

	//
	if settings.AllowIps, err = settings.NewAllowedKeys(config, "allow_ips"); err != nil {
		return err
	}

	if settings.AllowApiKeys, err = settings.NewAllowedKeys(config, "allow_api_keys"); err != nil {
		return err
	}

	if settings.AllowApiKeys.Enable && !_Tls.Enable {
		msg := "enabled api keys without using tls"
		fmt.Printf("!!! WARNING %s\n", msg)
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

	engine.NoRoute(func(ctx *gin.Context) {
		// ctx.Redirect(http.StatusFound, "https://example.local/" + c.Request.URL.Path)
		// ctx.String(http.StatusNotFound, "not found")
		settings.ReqLogger.Warn(
			"route not found",
			zap.String("ip", ctx.ClientIP()),
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()),
		)

		ctx.JSON(http.StatusNotFound, gin.H{"code": -1, "msg": "route not found"})
	})

	router = &engine.RouterGroup
	router.Static("/site", "./site")

	// TODO: apply aipLogger
	handlers.RouteOpen(router)
	handlers.RouteChatgpt(router, auth)
	handlers.RouteAuth(router, auth)

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

func logFields(ctx *gin.Context, start time.Time, lens ...uint) (fields []zap.Field) {
	if len(lens) > 0 {
		fields = make([]zap.Field, 0, lens[0])
	} else {
		fields = make([]zap.Field, 0, 6)
	}

	appendString := func(key, val string) {
		fields = append(fields, zap.String(key, val))
	}

	// requestId := uuid.NewString()
	appendString("ip", ctx.ClientIP())
	appendString("method", ctx.Request.Method)
	appendString("path", ctx.Request.URL.Path)
	appendString("query", ctx.Request.URL.RawQuery)

	status := ctx.Writer.Status()
	latencyMs := float64(time.Since(start).Microseconds()) / 1e3

	fields = append(fields, zap.Int("status", status))
	fields = append(fields, zap.Float64("latencyMs", latencyMs))

	return fields
}

func aipLogger(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()

	fields := logFields(ctx, start)

	status := ctx.Writer.Status()
	switch {
	case status < 400:
		settings.ReqLogger.Info("api request", fields...)
	case status < 500:
		settings.ReqLogger.Warn("api request", fields...)
	default:
		settings.ReqLogger.Error("api request", fields...)
	}
}
