package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	*gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db}
}
