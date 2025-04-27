package dtos

import (
	"Common"
	"Crypto_Bot/MainServer/github_sdk"
)

func ConvertCommit(commit *github_sdk.Commit, chatId int64, link string) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    chatId,
		Link:      link,
		Event:     Common.Commit,
		Author:    commit.Commit.Author.Name,
		Title:     commit.Commit.Message,
		UpdatedAt: commit.Commit.Committer.Date,
	}
}

func ConvertIssue(issue *github_sdk.Issue, chatId int64, link string) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    chatId,
		Link:      link,
		Event:     Common.Issue,
		Author:    issue.User.Login,
		Title:     issue.Title,
		UpdatedAt: issue.UpdatedAt,
	}
}

func ConvertPR(pr *github_sdk.PullRequest, chatId int64, link string) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    chatId,
		Link:      link,
		Event:     Common.PullRequest,
		Author:    pr.User.Login,
		Title:     pr.Title,
		UpdatedAt: pr.UpdatedAt,
	}
}
