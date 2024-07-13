package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spy-cat-api/env"

	"spy-cat-api/middleware/logger"
)

type Handler struct {
	env *env.Environment
}

func NewHandler(env *env.Environment) *Handler {
	return &Handler{env: env}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		logger.RequestLogger(),
		logger.ResponseLogger(),
	)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		h.initCatRoutes(api)
		h.initMissionRoutes(api)
	}
}
