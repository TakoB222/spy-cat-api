package handlers

import (
	"errors"
	"fmt"
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
		api.PUT("/:id/assign", h.AssignToMission)
		api.PUT("/targets/:target_id", h.UpdateTarget)
		api.DELETE("/targets/:target_id", h.DeleteTarget)
	}
}

func (h *Handler) CreateMission(c *gin.Context) {
	var mission models.Mission
	if err := c.ShouldBindJSON(&mission); err != nil {
		helper.Error(c, http.StatusBadRequest, fmt.Errorf("failed unmarshal request: %v", err))
		return
	}
	if err := mission.Validate(); err != nil {
		helper.Error(c, http.StatusBadRequest, err)
		return
	}

	h.env.Storage.Create(&mission)
	helper.Ok(c, mission)
}

func (h *Handler) GetMissions(c *gin.Context) {
	var missions []models.Mission
	h.env.Storage.Preload("Targets").Find(&missions)
	helper.Ok(c, missions)
}

func (h *Handler) GetMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := h.env.Storage.Preload("Targets").First(&mission, id).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, nil)
			return
		}
		helper.CriticalError(c, err)
		return
	}
	helper.Ok(c, mission)
}

func (h *Handler) UpdateMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := h.env.Storage.First(&mission, id).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, nil)
			return
		}
		helper.CriticalError(c, err)
		return
	}

	var input types.UpdateMissionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, fmt.Errorf("failed unmarshal request: %v", err))
		return
	}

	if mission.Completed {
		helper.Error(c, http.StatusBadRequest, errors.New("mission already completed"))
		return
	}

	h.env.Storage.Model(&mission).Updates(input)
	helper.Ok(c, mission)
}

func (h *Handler) DeleteMission(c *gin.Context) {
	id := c.Param("id")
	var mission models.Mission
	if err := h.env.Storage.First(&mission, id).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, nil)
			return
		}
		helper.CriticalError(c, err)
		return
	}

	if mission.CatID != 0 {
		helper.Error(c, http.StatusBadRequest, errors.New("cannot delete assigned mission"))
		return
	}

	h.env.Storage.Delete(&mission)
	helper.Ok(c, nil)
}

func (h *Handler) AddTarget(c *gin.Context) {
	missionID := c.Param("mission_id")
	var mission models.Mission
	if err := h.env.Storage.First(&mission, missionID).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, errors.New("mission not found"))
			return
		}
		helper.CriticalError(c, err)
		return
	}

	if mission.Completed {
		helper.Error(c, http.StatusBadRequest, errors.New("cannot add target to completed mission"))
		return
	}
	if len(mission.Targets) >= 3 {
		helper.Error(c, http.StatusBadRequest, errors.New("max targets per mission - 3"))
		return
	}

	var target models.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		helper.Error(c, http.StatusBadRequest, fmt.Errorf("failed unmarshal request: %v", err))
		return
	}

	target.MissionID = mission.ID
	h.env.Storage.Create(&target)
	helper.Ok(c, target)
}

func (h *Handler) AssignToMission(c *gin.Context) {
	missionID := c.Param("id")
	var mission models.Mission
	if err := h.env.Storage.First(&mission, missionID).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, errors.New("mission not found"))
			return
		}
		helper.CriticalError(c, err)
		return
	}

	if mission.Completed {
		helper.Error(c, http.StatusBadRequest, errors.New("cannot assign to completed mission"))
		return
	}
	if mission.CatID != 0 {
		helper.Error(c, http.StatusBadRequest, errors.New("mission already assigned"))
		return
	}

	var input types.AssignToMissionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, fmt.Errorf("failed unmarshal request: %v", err))
		return
	}

	currentMission := models.Mission{}
	if err := h.env.Storage.First(&currentMission, map[string]interface{}{
		"cat_id": input.CatID,
	}).Error; err != nil {
		if !models.IsErrorNotFound(err) {
			helper.CriticalError(c, err)
			return
		}
	}
	if currentMission.ID != 0 {
		helper.Error(c, http.StatusBadRequest, errors.New("cat already assigned to another mission"))
		return
	}

	mission.CatID = input.CatID
	h.env.Storage.Updates(mission)
	helper.Ok(c, nil)
}

func (h *Handler) UpdateTarget(c *gin.Context) {
	targetID := c.Param("target_id")
	var target models.Target
	if err := h.env.Storage.First(&target, targetID).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, nil)
			return
		}
		helper.CriticalError(c, err)
		return
	}

	var mission models.Mission
	if err := h.env.Storage.First(&mission, target.MissionID).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, errors.New("mission not found"))
			return
		}
		helper.CriticalError(c, err)
		return
	}

	if mission.Completed || target.Completed {
		helper.Error(c, http.StatusBadRequest, errors.New("cannot update completed mission/target"))
		return
	}

	var input types.UpdateTargetRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, fmt.Errorf("failed unmarshal request: %v", err))
		return
	}

	h.env.Storage.Model(&target).Updates(input)
	helper.Ok(c, target)
}

func (h *Handler) DeleteTarget(c *gin.Context) {
	targetID := c.Param("target_id")
	var target models.Target
	if err := h.env.Storage.First(&target, targetID).Error; err != nil {
		if models.IsErrorNotFound(err) {
			helper.Error(c, http.StatusNotFound, nil)
			return
		}
		helper.CriticalError(c, err)
		return
	}

	if target.Completed {
		helper.Error(c, http.StatusBadRequest, errors.New("cannot delete completed target"))
	}
}
