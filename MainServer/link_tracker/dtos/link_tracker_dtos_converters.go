package dtos

import (
	"Common"
	"Crypto_Bot/MainServer/github_sdk"
)

func ConvertCommit(commit *github_sdk.Commit, chatId int64) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId: chatId,
		Link:   commit.URL,
		Event:  Common.Commit,
		Author: commit.Commit.Author.Name,
		Title:  commit.Commit.Message,
	}
}

func ConvertIssue(issue *github_sdk.Issue, chatId int64) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    chatId,
		Link:      issue.URL,
		Event:     Common.Issue,
		Author:    issue.User.Login,
		Title:     issue.Title,
		UpdatedAt: issue.UpdatedAt,
	}
}

func ConvertPR(pr *github_sdk.PullRequest, chatId int64) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    chatId,
		Link:      pr.URL,
		Event:     Common.PullRequest,
		Author:    pr.User.Login,
		Title:     pr.Title,
		UpdatedAt: pr.UpdatedAt,
	}
}
