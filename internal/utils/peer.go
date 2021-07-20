package utils

import "github.com/gotd/td/tg"

func GetPeerID(peer tg.PeerClass) int {
	switch peer := peer.(type) {
	case *tg.PeerChat:
		return peer.ChatID
	case *tg.PeerChannel:
		return peer.ChannelID
	case *tg.PeerUser:
		return peer.UserID
	default:
		return 0
	}
}

func GetInputPeerID(peer tg.InputPeerClass) int {
	switch peer := peer.(type) {
	case *tg.InputPeerChannel:
		return peer.ChannelID
	case *tg.InputPeerChat:
		return peer.ChatID
	case *tg.InputPeerUser:
		return peer.UserID
	default:
		return 0
	}
}
