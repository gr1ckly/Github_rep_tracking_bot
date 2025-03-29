package storage

import "time"

const (
	Commit      string = "commit"
	Issue       string = "issue"
	PullRequest string = "pull_request"
)

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
	ChatID int
	Type   string
}

type ChatRepoRecord struct {
	ID     int
	Chat   *Chat
	Repo   *Repo
	Tags   []string
	Events []string
}
