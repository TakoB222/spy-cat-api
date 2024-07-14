package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"spy-cat-api/handlers/helper"
	"spy-cat-api/handlers/types"
	"spy-cat-api/models"
	"spy-cat-api/utils"
)

func (h *Handler) initCatRoutes(group *gin.RouterGroup) {
	api := group.Group("/cats")
	{
		api.POST("", h.CreateCat)
		api.GET("", h.GetCats)
		api.GET("/:id", h.GetCat)
		api.PUT("/:id", h.UpdateCat)
		api.DELETE("/:id", h.DeleteCat)
	}
}

func (h *Handler) CreateCat(c *gin.Context) {
	var cat models.Cat
	if err := c.ShouldBindJSON(&cat); err != nil {
		helper.Error(c, http.StatusBadRequest, fmt.Errorf("failed unmarshal request: %v", err))
		return
	}

	if !utils.ValidateCatBreed(cat.Breed) {
		helper.Error(c, http.StatusBadRequest, errors.New("invalid cat breed"))
		return
	}

	if err := h.env.Storage.Create(&cat).Error; err != nil {
		helper.CriticalError(c, err)
		return
	}
	helper.Ok(c, cat)
}

func (h *Handler) GetCats(c *gin.Context) {
	var cats []models.Cat
	h.env.Storage.Find(&cats)
	helper.Ok(c, cats)
}

func (h *Handler) GetCat(c *gin.Context) {
	id := c.Param("id")
	var cat models.Cat
	if err := h.env.Storage.First(&cat, id).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, nil)
			return
		}
		helper.CriticalError(c, err)
		return
	}
	helper.Ok(c, cat)
}

func (h *Handler) UpdateCat(c *gin.Context) {
	id := c.Param("id")
	var cat models.Cat
	if err := h.env.Storage.First(&cat, id).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, nil)
			return
		}
		helper.CriticalError(c, err)
		return
	}

	var input types.UpdateCatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, fmt.Errorf("failed unmarshal request: %v", err))
		return
	}
	if err := input.Validate(); err != nil {
		helper.Error(c, http.StatusBadRequest, err)
		return
	}

	h.env.Storage.Model(&cat).Updates(input)
	helper.Ok(c, cat)
}

func (h *Handler) DeleteCat(c *gin.Context) {
	id := c.Param("id")
	if err := h.env.Storage.Delete(&models.Cat{}, id).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, nil)
			return
		}
		helper.CriticalError(c, err)
		return
	}
	helper.Ok(c, nil)
}
