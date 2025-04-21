package entities

import "time"

type Message struct {
	UserId    int64     `json:"user_id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}
