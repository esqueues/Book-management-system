package models

import "time"

type Chat struct {
	ID       string
	ClientID string
	AdminID  string
	Messages []Message
	Active   bool
}

type Message struct {
	SenderID  string
	Content   string
	Timestamp time.Time
}
