package env

import (
	"errors"
	"spy-cat-api/models"
)

var NilConfigErr = errors.New("no config")

type Environment struct {
	Storage *models.Storage
	config  *Config
}
