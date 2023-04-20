package handlers

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/d2jvkpn/x-ai/pkg/chatgpt"
	"github.com/gin-gonic/gin"
)

func chatCompl(ctx *gin.Context) {
	var (
		err  error
		code int
		req  *chatgpt.CompReq
		res  *chatgpt.CompRes
	)

	defer func() {
		if err != nil {
			log.Printf("!!! completions: code: %d, error: %v\n", code, err)
		}
	}()

	req = new(chatgpt.CompReq)
	if err = ctx.BindJSON(req); err != nil {
		code = -1
		ctx.JSON(http.StatusBadRequest, gin.H{"code": code, "msg": "unmarshal failed"})
		return
	}
	if err = req.Validate(); err != nil {
		code = -2
		ctx.JSON(http.StatusBadRequest, gin.H{"code": code, "msg": err.Error()})
		return
	}

	err_msg := "failed to call third party services"
	if res, err = settings.GPTCli.Completions(ctx, req); err != nil {
		code = 11
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"code": code, "msg": err_msg})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func imgGen(ctx *gin.Context) {
	var (
		err  error
		code int
		req  *chatgpt.ImgGenReq
		res  *chatgpt.ImgGenRes
	)

	defer func() {
		if err != nil {
			log.Printf("!!! imggen: code: %d, error: %v\n", code, err)
		}
	}()

	req = new(chatgpt.ImgGenReq)
	if err = ctx.BindJSON(&req); err != nil {
		code = -1
		ctx.JSON(http.StatusBadRequest, gin.H{"code": code, "msg": "unmarshal failed"})
		return
	}
	if err = req.Validate(); err != nil {
		code = -2
		ctx.JSON(http.StatusBadRequest, gin.H{"code": code, "msg": err.Error()})
		return
	}

	err_msg := "failed to call third party services"
	if res, err = settings.GPTCli.ImgGen(ctx, req); err != nil {
		code = 11
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"code": code, "msg": err_msg})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func models(ctx *gin.Context) {
	var (
		err  error
		code int
		res  *chatgpt.ModelsRes
	)

	defer func() {
		if err != nil {
			log.Printf("!!! models: code: %d, error: %v\n", code, err)
		}
	}()

	err_msg := "failed to call third party services"
	if res, err = settings.GPTCli.Models(ctx); err != nil {
		code = 11
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"code": code, "msg": err_msg})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
