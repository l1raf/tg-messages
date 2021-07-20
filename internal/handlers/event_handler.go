package handlers

import (
	"context"
	"log"
	"tg-messages/internal/database"
	"tg-messages/internal/models"
	"tg-messages/internal/msg"
	"tg-messages/internal/utils"

	"github.com/gotd/td/tg"
)

type MessagesHandler struct {
	msgRepo *database.MessageRepository
	chats   []int
}

func NewMessagesHandler(ms *database.MessageRepository, chats []int) *MessagesHandler {
	return &MessagesHandler{
		msgRepo: ms,
		chats:   chats,
	}
}

func (mh *MessagesHandler) EditMessageHandler(ctx context.Context, entities tg.Entities, u *tg.UpdateEditMessage) error {
	message, ok := u.Message.(*tg.Message)

	if !ok || message.Out {
		return nil
	}

	if !msg.Contains(mh.chats, utils.GetPeerID(message.PeerID)) {
		return nil
	}

	err := mh.msgRepo.Update(models.Message{
		MessageID:     message.ID,
		FromID:        utils.GetPeerID(message.FromID),
		PeerID:        utils.GetPeerID(message.PeerID),
		Text:          message.Message,
		Date:          message.Date,
		EditDate:      message.EditDate,
		PostAuthor:    message.PostAuthor,
		Outgoing:      message.Out,
		Mentioned:     message.Mentioned,
		Silent:        message.Silent,
		Post:          message.Post,
		FromScheduled: message.FromScheduled,
		Pinned:        message.Pinned,
		ViaBotID:      message.ViaBotID,
		Views:         message.Views,
		Forwards:      message.Forwards,
	})

	log.Println("OnEditMessage:", message)

	return err
}

func (mh *MessagesHandler) NewMessageHandler(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
	message, ok := u.Message.(*tg.Message)

	if !ok || message.Out {
		return nil
	}

	err := handleMessage(message, mh)
	log.Println("OnNewMessage:", message)

	return err
}

func (mh *MessagesHandler) NewChannelMessageHandler(ctx context.Context, entities tg.Entities, u *tg.UpdateNewChannelMessage) error {
	message, ok := u.Message.(*tg.Message)

	if !ok || message.Out {
		return nil
	}

	err := handleMessage(message, mh)
	log.Println("OnNewChannelMessage:", message)

	return err
}

func handleMessage(message *tg.Message, mh *MessagesHandler) error {
	if !msg.Contains(mh.chats, utils.GetPeerID(message.PeerID)) {
		return nil
	}

	return mh.msgRepo.Create(models.Message{
		MessageID:     message.ID,
		FromID:        utils.GetPeerID(message.FromID),
		PeerID:        utils.GetPeerID(message.PeerID),
		Text:          message.Message,
		Date:          message.Date,
		EditDate:      message.EditDate,
		PostAuthor:    message.PostAuthor,
		Outgoing:      message.Out,
		Mentioned:     message.Mentioned,
		Silent:        message.Silent,
		Post:          message.Post,
		FromScheduled: message.FromScheduled,
		Pinned:        message.Pinned,
		ViaBotID:      message.ViaBotID,
		Views:         message.Views,
		Forwards:      message.Forwards,
	})
}
