package msg

import (
	"context"
	"log"

	"tg-messages/internal/utils"

	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/telegram/query/dialogs"
	"github.com/gotd/td/telegram/query/messages"
	"github.com/gotd/td/tg"
)

func GetLastNMessages(ctx context.Context, api *tg.Client, chats []int, numOfMessagesToGet int) ([]tg.Message, error) {
	var tgMessages []tg.Message

	cb := func(ctx context.Context, dlg dialogs.Elem) error {
		// Skip deleted dialogs
		if _, empty := dlg.Peer.(*tg.InputPeerEmpty); empty {
			return nil
		}

		if !Contains(chats, utils.GetInputPeerID(dlg.Peer)) {
			return nil
		}

		i := 0
		count := numOfMessagesToGet

		f := func(elem messages.Elem) error {
			msg, ok := elem.Msg.(*tg.Message)

			if !ok {
				count++
				return nil
			}

			i++
			tgMessages = append(tgMessages, *msg)
			log.Printf("%d. %v: message: %v", i, msg.PeerID, msg.Message)

			return nil
		}

		reverse(tgMessages)
		iter := dlg.Messages(api).BatchSize(100).Iter()

		for i := 0; i < count && iter.Next(ctx); i++ {
			if err := f(iter.Value()); err != nil {
				return err
			}
		}

		return iter.Err()
	}

	err := query.GetDialogs(api).ForEach(ctx, cb)

	return tgMessages, err
}

func reverse(messages []tg.Message) {
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
}

func Contains(chats []int, val int) bool {
	for i := range chats {
		if chats[i] == val {
			return true
		}
	}

	return false
}
