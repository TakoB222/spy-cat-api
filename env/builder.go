package env

import (
	"errors"
	"spy-cat-api/models"
	"spy-cat-api/services"
)

var BuiltAlreadyErr = errors.New("already built")

type Builder struct {
	env *Environment
	err error
}

func NewBuilder() *Builder {
	return &Builder{
		env: &Environment{},
	}
}

func (b *Builder) SetConfig(c *Config) *Builder {
	if b.err != nil {
		return b
	}

	if c == nil {
		b.err = errors.New("cannot set nil config")
		return b
	}

	b.env.config = c
	return b
}

func (b *Builder) ConnectToPostgresDB() *Builder {
	if b.err != nil {
		return b
	}

	if b.env.config == nil {
		b.err = NilConfigErr
		return b
	}

	db, err := services.NewDB(b.env.config.DB)
	if err != nil {
		b.err = err
		return b
	}

	b.env.Storage = models.NewStorage(db)
	return b
}

func (b *Builder) Build() (*Environment, error) {
	if b.err != nil {
		return nil, b.err
	}

	b.err = BuiltAlreadyErr
	return b.env, nil
}
