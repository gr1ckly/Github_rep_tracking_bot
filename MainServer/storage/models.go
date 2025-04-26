package storage

import "time"

type Repo struct {
	ID         int
	Name       string
	Owner      string
	Link       string
	LastCommit time.Time
	LastIssue  time.Time
	LastPR     time.Time
}

type Chat struct {
	ChatID int64
	Type   string
}

type ChatRepoRecord struct {
	ID     int
	Chat   *Chat
	Repo   *Repo
	Tags   []string
	Events []string
}
