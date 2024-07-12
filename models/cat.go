package models

import (
	"gorm.io/gorm"
)

type Cat struct {
	gorm.Model
	Name              string `json:"name" binding:"required"`
	YearsOfExperience int64  `json:"years_of_experience" binding:"required"`
	Breed             string `json:"breed" binding:"required"`
	Salary            int64  `json:"salary" binding:"required"`
}
