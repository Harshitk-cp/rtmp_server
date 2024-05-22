package room

import "time"

type ChatMessage struct {
	SenderID  string
	RoomID    string
	Content   string
	Timestamp time.Time
}
