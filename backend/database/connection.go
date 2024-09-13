package database

import (
	"backend/database/drivers"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	config := &gorm.Config{}

	// Engine
	switch viper.GetString("database.driver") {
	case "postgres":
		db, err = drivers.ConnectPostgres(config)
	default:
		panic("Invalid database driver")
	}
	DB = db

	return db, err
}
