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
	msgStore *database.MessageRepository
	chats    []int
}

func NewMessagesHandler(ms *database.MessageRepository, chats []int) *MessagesHandler {
	return &MessagesHandler{
		msgStore: ms,
		chats:    chats,
	}
}

func (mh *MessagesHandler) EditMessageHandler(ctx context.Context, entities tg.Entities, u *tg.UpdateEditMessage) error {
	message, ok := u.Message.(*tg.Message)

	if !ok || message.Out {
		return nil
	}

	if !msg.Contains(mh.chats, utils.GetPeerId(message.PeerID)) {
		return nil
	}

	err := mh.msgStore.Update(models.Message{
		MessageId:     message.ID,
		FromId:        utils.GetPeerId(message.FromID),
		PeerId:        utils.GetPeerId(message.PeerID),
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
	if !msg.Contains(mh.chats, utils.GetPeerId(message.PeerID)) {
		return nil
	}

	return mh.msgStore.Create(models.Message{
		MessageId:     message.ID,
		FromId:        utils.GetPeerId(message.FromID),
		PeerId:        utils.GetPeerId(message.PeerID),
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
	})
}
