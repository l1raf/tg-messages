package database

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"

	"tg-messages/internal/models"
)

var db *gorm.DB

func OpenConnection(connectionString string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open("postgres", connectionString)
	return db, err
}

func Init() error {
	return db.AutoMigrate(&models.Message{}).Error
}
