package github_sdk

import (
	"Crypto_Bot/MainServer/server/validators"
	"time"
)

type GithubService interface {
	GetCommits(repoName string, owner string, since time.Time) ([]Commit, error)
	GetIssues(repoName string, owner string, since time.Time) ([]Issue, error)
	GetPullRequests(repoName string, owner string, since time.Time) ([]PullRequest, error)
	validators.Checker[string, bool]
}
