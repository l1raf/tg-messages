package database

import (
	"tg-messages/internal/models"

	"github.com/jinzhu/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageStore(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (mr *MessageRepository) Create(msg models.Message) error {
	return db.Create(&msg).Error
}

func (mr *MessageRepository) Update(msg models.Message) error {
	var foundMsg models.Message

	err := db.
		Where("message_id = ?", msg.MessageId).
		Where("peer_id = ?", msg.PeerId).
		Find(&foundMsg).Error

	if err != nil {
		return err
	}

	return db.Model(&foundMsg).Update(&msg).Error
}

func (mr *MessageRepository) GetAll() ([]models.Message, error) {
	var messages []models.Message

	err := db.Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
