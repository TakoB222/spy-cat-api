package helper

import (
	"github.com/gin-gonic/gin"

	"spy-cat-api/pkg/logger"
)

type response struct {
	Message string `json:"message"`
}

func NewResponse(ctx *gin.Context, statusCode int, message string) {
	logger.Error(message)
	ctx.AbortWithStatusJSON(statusCode, response{message})
}
