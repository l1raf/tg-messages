package handlers

import (
	"context"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
)

var _ updates.Handler = (*GapAdapter)(nil)

type GapAdapter struct {
	next telegram.UpdateHandler
	ctx  context.Context
}

func NewGapAdapter(ctx context.Context, next telegram.UpdateHandler) *GapAdapter {
	return &GapAdapter{next, ctx}
}

func (a *GapAdapter) HandleDiff(diff updates.DiffUpdate) error {
	return a.next.Handle(a.ctx, &tg.Updates{
		Updates: append(
			msgsToUpdates(diff.NewMessages),
			encryptedMsgsToUpdates(diff.NewEncryptedMessages)...,
		),
		Users: diff.Users,
		Chats: diff.Chats,
	})
}

func (a *GapAdapter) HandleUpdates(ents *updates.Entities, ups []tg.UpdateClass) error {
	return a.next.Handle(a.ctx, &tg.Updates{
		Updates: ups,
		Users:   ents.AsUsers(),
		Chats:   ents.AsChats(),
	})
}

func (a *GapAdapter) ChannelTooLong(channelID int) {
	log.Print("ChannelTooLong:", channelID)
}

func msgsToUpdates(msgs []tg.MessageClass) []tg.UpdateClass {
	var updates []tg.UpdateClass
	for _, msg := range msgs {
		updates = append(updates, &tg.UpdateNewMessage{
			Message:  msg,
			Pts:      -1,
			PtsCount: -1,
		})
	}
	return updates
}

func encryptedMsgsToUpdates(msgs []tg.EncryptedMessageClass) []tg.UpdateClass {
	var updates []tg.UpdateClass
	for _, msg := range msgs {
		updates = append(updates, &tg.UpdateNewEncryptedMessage{
			Message: msg,
			Qts:     -1,
		})
	}
	return updates
}
