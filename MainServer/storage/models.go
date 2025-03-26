package storage

import "time"

type event string

type Tag string

const (
	Commit      event = "commit"
	Issue       event = "issue"
	PullRequest event = "pull_request"
)

type Repo struct {
	ID         int
	Name       string
	Owner      string
	LastCommit time.Time
	LastIssue  time.Time
	LastPR     time.Time
}

type Chat struct {
	ID     int
	ChatID int
	Type   string
}

type ChatRepoRecord struct {
	ID     int
	Chat   Chat
	Repo   Repo
	Tags   []Tag
	Events []event
}
