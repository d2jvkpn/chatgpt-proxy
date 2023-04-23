package handlers

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/d2jvkpn/chatgpt-proxy/internal/settings"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func listModels(ctx *gin.Context) {
	var (
		err      error
		apiError *openai.APIError
		res      openai.ModelsList
	)

	defer func() {
		if err != nil {
			log.Printf("!!! listModels error: %+v\n", apiError)
		}
	}()

	if res, err = settings.GPTCli2.ListModels(ctx); err != nil {
		apiError, _ = HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func createCompletions(ctx *gin.Context) {
	var (
		err      error
		apiError *openai.APIError
		req      openai.CompletionRequest
		res      openai.CompletionResponse
	)

	defer func() {
		if err != nil {
			log.Printf("!!! createCompletions error: %+v\n", apiError)
		}
	}()

	if err = ctx.BindJSON(&req); err != nil {
		apiError = unmarshalErr
		ctx.JSON(apiError.HTTPStatusCode, apiError)
		return
	}

	if res, err = settings.GPTCli2.CreateCompletion(ctx, req); err != nil {
		apiError, _ = HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func createChatCompletions(ctx *gin.Context) {
	var (
		err      error
		apiError *openai.APIError
		req      openai.ChatCompletionRequest
		res      openai.ChatCompletionResponse
	)

	defer func() {
		if err != nil {
			log.Printf("!!! createChatCompletions error: %+v\n", apiError)
		}
	}()

	if err = ctx.BindJSON(&req); err != nil {
		apiError = unmarshalErr
		ctx.JSON(apiError.HTTPStatusCode, apiError)
		return
	}

	if res, err = settings.GPTCli2.CreateChatCompletion(ctx, req); err != nil {
		apiError, _ = HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func createImage(ctx *gin.Context) {
	var (
		err      error
		apiError *openai.APIError
		req      openai.ImageRequest
		res      openai.ImageResponse
	)

	defer func() {
		if err != nil {
			log.Printf("!!! createImage error: %+v\n", apiError)
		}
	}()

	if err = ctx.BindJSON(&req); err != nil {
		apiError = unmarshalErr
		ctx.JSON(apiError.HTTPStatusCode, apiError)
		return
	}

	if res, err = settings.GPTCli2.CreateImage(ctx, req); err != nil {
		apiError, _ = HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func createEditImage(ctx *gin.Context) {
	var (
		err      error
		apiError *openai.APIError
		req      openai.ImageEditRequest
		res      openai.ImageResponse
	)

	defer func() {
		if err != nil {
			log.Printf("!!! createEditImage error: %+v\n", apiError)
		}
	}()

	if err = ctx.BindJSON(&req); err != nil {
		apiError = unmarshalErr
		ctx.JSON(apiError.HTTPStatusCode, apiError)
		return
	}

	if res, err = settings.GPTCli2.CreateEditImage(ctx, req); err != nil {
		apiError, _ = HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func createVariImage(ctx *gin.Context) {
	var (
		err      error
		apiError *openai.APIError
		req      openai.ImageVariRequest
		res      openai.ImageResponse
	)

	defer func() {
		if err != nil {
			log.Printf("!!! createVariImage error: %+v\n", apiError)
		}
	}()

	if err = ctx.BindJSON(&req); err != nil {
		apiError = unmarshalErr
		ctx.JSON(apiError.HTTPStatusCode, apiError)
		return
	}

	if res, err = settings.GPTCli2.CreateVariImage(ctx, req); err != nil {
		apiError, _ = HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func createEmbeddings(ctx *gin.Context) {
	var (
		err      error
		apiError *openai.APIError
		req      openai.EmbeddingRequest
		res      openai.EmbeddingResponse
	)

	defer func() {
		if err != nil {
			log.Printf("!!! createEmbeddings error: %+v\n", apiError)
		}
	}()

	if err = ctx.BindJSON(&req); err != nil {
		apiError = unmarshalErr
		ctx.JSON(apiError.HTTPStatusCode, apiError)
		return
	}

	if res, err = settings.GPTCli2.CreateEmbeddings(ctx, req); err != nil {
		apiError, _ = HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func listEngines(ctx *gin.Context) {
	var (
		err      error
		apiError *openai.APIError
		res      openai.EnginesList
	)

	defer func() {
		if err != nil {
			log.Printf("!!! listEngines error: %+v\n", apiError)
		}
	}()

	if res, err = settings.GPTCli2.ListEngines(ctx); err != nil {
		apiError, _ = HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
