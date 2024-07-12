package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spy-cat-api/handlers/helper"
	"spy-cat-api/handlers/types"
	"spy-cat-api/models"
)

func (h *Handler) initMissionRoutes(group *gin.RouterGroup) {
	api := group.Group("/missions")
	{
		api.POST("", h.CreateMission)
		api.GET("", h.GetMissions)
		api.GET("/:id", h.GetMission)
		api.PUT("/:id", h.UpdateMission)
		api.DELETE("/:id", h.DeleteMission)
		api.POST("/:mission_id/targets", h.AddTarget)
		api.PUT("/targets/:target_id", h.UpdateTarget)
		api.DELETE("/targets/:target_id", h.DeleteTarget)
	}
}

func (h *Handler) CreateMission(c *gin.Context) {
	var mission models.Mission
	if err := c.ShouldBindJSON(&mission); err != nil {
		helper.NewResponse(c, http.StatusBadRequest, "filed unmarshal request")
		return
	}
	if err := mission.Validate(); err != nil {
		helper.NewResponse(c, http.StatusBadRequest, err.Error())
	}

	h.env.Storage.Create(&mission)
	c.JSON(http.StatusOK, mission)
}

func (h *Handler) GetMissions(c *gin.Context) {
	var missions []models.Mission
	h.env.Storage.Preload("Targets").Find(&missions)
	c.JSON(http.StatusOK, missions)
}

func (h *Handler) GetMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := h.env.Storage.Preload("Targets").First(&mission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
	}
	c.JSON(http.StatusOK, mission)
}

func (h *Handler) UpdateMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := h.env.Storage.First(&mission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	var input types.UpdateMissionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.NewResponse(c, http.StatusBadRequest, "filed unmarshal request")
		return
	}

	if mission.Completed {
		helper.NewResponse(c, http.StatusBadRequest, "mission already completed")
		return
	}

	h.env.Storage.Model(&mission).Updates(input)
	c.JSON(http.StatusOK, mission)
}

func (h *Handler) DeleteMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := h.env.Storage.First(&mission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	if mission.CatID != 0 {
		helper.NewResponse(c, http.StatusBadRequest, "cannot delete assigned mission")
		return
	}

	h.env.Storage.Delete(&mission)
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) AddTarget(c *gin.Context) {
	missionID := c.Param("mission_id")
	var mission models.Mission
	if err := h.env.Storage.First(&mission, missionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	if mission.Completed {
		helper.NewResponse(c, http.StatusBadRequest, "cannot add target to completed mission")
		return
	}

	var target models.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		helper.NewResponse(c, http.StatusBadRequest, "filed unmarshal request")
		return
	}

	target.MissionID = mission.ID
	h.env.Storage.Create(&target)
	c.JSON(http.StatusOK, target)
}

func (h *Handler) UpdateTarget(c *gin.Context) {
	targetID := c.Param("target_id")
	var target models.Target
	if err := h.env.Storage.First(&target, targetID).Error; err != nil {
		helper.NewResponse(c, http.StatusNotFound, "target not found")
		return
	}

	var mission models.Mission
	if err := h.env.Storage.First(&mission, target.MissionID).Error; err != nil {
		helper.NewResponse(c, http.StatusNotFound, "mission not found")
		return
	}

	if mission.Completed || target.Completed {
		helper.NewResponse(c, http.StatusBadRequest, "cannot update completed mission/target")
		return
	}

	var input types.UpdateTargetRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.env.Storage.Model(&target).Updates(input)
	c.JSON(http.StatusOK, target)
}

func (h *Handler) DeleteTarget(c *gin.Context) {
	targetID := c.Param("target_id")
	var target models.Target
	if err := h.env.Storage.First(&target, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	if target.Completed {
		helper.NewResponse(c, http.StatusBadRequest, "cannot delete completed target")
		return
	}
}
