package dtos

import (
	"Common"
	"Crypto_Bot/MainServer/github_sdk"
	"Crypto_Bot/MainServer/storage"
)

func ConvertCommit(commit *github_sdk.Commit, record *storage.ChatRepoRecord) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    record.Chat.ChatID,
		Link:      record.Repo.Link,
		Event:     Common.Commit,
		Author:    commit.Commit.Author.Name,
		Title:     commit.Commit.Message,
		UpdatedAt: commit.Committer.Date,
	}
}

func ConvertIssue(issue *github_sdk.Issue, record *storage.ChatRepoRecord) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    record.Chat.ChatID,
		Link:      record.Repo.Link,
		Event:     Common.Issue,
		Author:    issue.User.Login,
		Title:     issue.Title,
		UpdatedAt: issue.UpdatedAt,
	}
}

func ConvertPR(pr *github_sdk.PullRequest, record *storage.ChatRepoRecord) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    record.Chat.ChatID,
		Link:      record.Repo.Link,
		Event:     Common.PullRequest,
		Author:    pr.User.Login,
		Title:     pr.Title,
		UpdatedAt: pr.UpdatedAt,
	}
}
