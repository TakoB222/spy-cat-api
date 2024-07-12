package models

import (
	"errors"
	"gorm.io/gorm"
)

type Mission struct {
	gorm.Model
	CatID     uint     `json:"cat_id" binding:"required" gorm:"constraint:OnDelete:CASCADE"`
	Completed bool     `json:"completed"`
	Targets   []Target `json:"targets" gorm:"constraint:OnDelete:CASCADE"`
}

func (m *Mission) Validate() error {
	if m.CatID == 0 {
		return errors.New("empty cat_id")
	}
	if len(m.Targets) > 3 {
		return errors.New("max targets per mission - 3")
	}

	return nil
}
