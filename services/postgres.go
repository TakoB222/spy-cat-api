package services

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"spy-cat-api/models"
)

type DBConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"port"`
}

func NewDB(c DBConfig) (db *gorm.DB, err error) {
	database, err := gorm.Open(postgres.Open(composeConnectionString(c)), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&models.Cat{}, &models.Mission{}, &models.Target{})

	return database, nil
}

func composeConnectionString(c DBConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		getTimezone(),
	)
}

func getTimezone() string {
	return "UTC"
}
