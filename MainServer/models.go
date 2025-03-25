package main

import "time"

type ChatId int64

type Link struct {
	Owner      string
	RepoName   string
	LastCommit time.Time
	LastIssue  time.Time
	LastPR     time.Time
}

type Tag string

type Event string

const (
	PullRequest Event = "pull_request"
	Commit      Event = "commit"
	Issue       Event = "issue"
)

type LinkCheck struct {
	ChatId ChatId
	Link   Link
	Tags   []Tag
	Events []Event
}
