package handlers

import (
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
		helper.NewResponse(c, http.StatusBadRequest, "filed unmarshal request")
		return
	}

	if !utils.ValidateCatBreed(cat.Breed) {
		helper.NewResponse(c, http.StatusBadRequest, "invalid cat breed")
		return
	}

	if err := h.env.Storage.Create(&cat).Error; err != nil {
		helper.NewResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) GetCats(c *gin.Context) {
	var cats []models.Cat
	h.env.Storage.Find(&cats)
	c.JSON(http.StatusOK, cats)
}

func (h *Handler) GetCat(c *gin.Context) {
	id := c.Param("id")
	var cat models.Cat
	if err := h.env.Storage.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) UpdateCat(c *gin.Context) {
	id := c.Param("id")
	var cat models.Cat
	if err := h.env.Storage.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
	}

	var input types.UpdateCatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.NewResponse(c, http.StatusBadRequest, "filed unmarshal request")
		return
	}
	if err := input.Validate(); err != nil {
		helper.NewResponse(c, http.StatusBadRequest, err.Error())
	}

	h.env.Storage.Model(&cat).Updates(input)
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) DeleteCat(c *gin.Context) {
	id := c.Param("id")
	if err := h.env.Storage.Delete(&models.Cat{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
	}
	c.JSON(http.StatusOK, nil)
}
