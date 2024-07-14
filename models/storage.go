package models

import (
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	*gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db}
}

func IsErrorNotFound(err error) bool {
	if !errors.As(err, &gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
