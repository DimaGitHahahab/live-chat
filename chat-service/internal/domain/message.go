package domain

import (
	"fmt"
	"time"
)

type Message struct {
	Id     int       `json:"id"`
	Sender string    `json:"sender"`
	Text   string    `json:"text"`
	SendAt time.Time `json:"send_at"`
}

func (m *Message) String() string {
	return fmt.Sprintf("%s: %s [%s]", m.Sender, m.Text, m.SendAt.Format("2 January 15:04"))
}
