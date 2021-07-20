package database

import (
	"tg-messages/internal/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func OpenConnection(connectionString string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open("postgres", connectionString)

	if err == nil {
		db.LogMode(false)
	}

	return db, err
}

func Init() error {
	return db.AutoMigrate(&models.Message{}).Error
}
