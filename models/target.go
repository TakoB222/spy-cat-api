package models

import "gorm.io/gorm"

type Target struct {
	gorm.Model
	MissionID uint   `json:"mission_id"`
	Name      string `json:"name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}
