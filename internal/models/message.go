package models

type Message struct {
	ID            int    `gorm:"primary_key"`
	MessageID     int    `gorm:"unique_index:idx_msg"` // ID of the message
	FromID        int    // ID of the sender of the message
	PeerID        int    `gorm:"unique_index:idx_msg"` // Peer ID, the chat where this message was sent
	Text          string // Message
	Date          int    // Date of the message (unix)
	EditDate      int    // Last edit date of this message
	PostAuthor    string // Name of the author of this message for channel posts (with signatures enabled)
	Outgoing      bool   // Is this an outgoing message
	Mentioned     bool   // Whether user was mentioned in this message
	Silent        bool   // Whether this is a silent message (no notification triggered)
	Post          bool   // Whether this is a channel post
	FromScheduled bool   // Whether this is a scheduled message
	Pinned        bool   // Whether this message is pinned
	ViaBotID      int    // ID of the inline bot that generated the message
	Views         int    // View count for channel posts
	Forwards      int    // Forward counter
}
