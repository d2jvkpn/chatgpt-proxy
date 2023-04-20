package handlers

import (
	"errors"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

var (
	unmarshalErr = &openai.APIError{
		Code: -100, Message: "unmarshal failed", Type: "unmarshal_error",
		HTTPStatusCode: http.StatusBadRequest,
	}
)

func HandleError(ctx *gin.Context, err error) (apiError *openai.APIError, ok bool) {
	var (
		theErr     error
		statusCode int
	)

	if theErr = errors.Unwrap(err); theErr == nil {
		statusCode = http.StatusInternalServerError
		apiError = &openai.APIError{
			Code: 100, Message: err.Error(), Type: "UNKNOWN",
			HTTPStatusCode: statusCode,
		}

		ctx.JSON(statusCode, apiError)
		return apiError, false
	}

	if apiError, ok = err.(*openai.APIError); !ok {
		statusCode = http.StatusInternalServerError
		apiError = &openai.APIError{
			Code: 101, Message: err.Error(), Type: "UNKNOWN",
			HTTPStatusCode: statusCode,
		}

		ctx.JSON(statusCode, apiError)
		return apiError, false
	}

	ctx.JSON(apiError.HTTPStatusCode, apiError)
	return apiError, true
}
