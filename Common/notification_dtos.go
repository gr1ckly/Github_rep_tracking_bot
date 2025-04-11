package Common

import "time"

type ChangingDTO struct {
	ChatId    int64     `json:"chat_id"`
	Link      string    `json:"link"`
	Event     string    `json:"event"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}
