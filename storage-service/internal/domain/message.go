package domain

import (
	"time"
)

type Message struct {
	Id     int       `json:"id"`
	Sender string    `json:"sender"`
	Text   string    `json:"text"`
	SendAt time.Time `json:"send_at"`
}
