package models

import (
	"errors"
	"gorm.io/gorm"
)

type Mission struct {
	gorm.Model
	CatID     uint     `json:"cat_id"`
	Completed bool     `json:"completed"`
	Targets   []Target `json:"targets"`
}

func (m *Mission) Validate() error {
	if len(m.Targets) > 3 {
		return errors.New("max targets per mission - 3")
	}
	return nil
}
