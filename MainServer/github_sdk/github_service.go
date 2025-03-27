package github_sdk

import "Crypto_Bot/MainServer/server/validators"

type GithubService interface {
	GetCommits(repoName string, owner string) ([]Commit, error)
	GetIssues(repoName string, owner string) ([]Issue, error)
	GetPullRequests(repoName string, owner string) ([]PullRequest, error)
	validators.Checker[string, bool]
}
