package main

type ChatDTO struct {
	ChatID int    `json:"chat_id"`
	Type   string `json:"type"`
}

type RepoDTO struct {
	Link   string   `json:"link"`
	ChatID int      `json:"chat_id"`
	Tags   []string `json:"tags"`
	Events []string `json:"events"`
}

type ErrorDTO struct {
	error string
}
