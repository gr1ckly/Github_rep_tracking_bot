package dtos

import "time"

type CommitChangesDTO struct {
	ChatId   int    `json:"chat_id"`
	Link     string `json:"link"`
	Author   string `json:"author"`
	Commiter string `json:"commiter"`
	Message  string `json:"message"`
	Branch   string `json:"branch"`
}

type IssueChangesSTO struct {
	ChatId    int       `json:"chat_id"`
	Link      string    `json:"link"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PullRequestChangesDTO struct {
	ChatId    int       `json:"chat_id"`
	Link      string    `json:"link"`
	Author    string    `json:"author"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}
