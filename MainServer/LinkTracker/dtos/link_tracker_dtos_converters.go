package dtos

import (
	"Crypto_Bot/MainServer/github_sdk"
	"Crypto_Bot/MainServer/storage"
)

func ConvertCommit(commit *github_sdk.Commit, chatId int) ChangingDTO {
	return ChangingDTO{
		ChatId: chatId,
		Link:   commit.URL,
		Event:  storage.Commit,
		Author: commit.Commit.Author.Name,
		Title:  commit.Commit.Message,
	}
}

func ConvertIssue(issue *github_sdk.Issue, chatId int) ChangingDTO {
	return ChangingDTO{
		ChatId:    chatId,
		Link:      issue.URL,
		Event:     storage.Issue,
		Author:    issue.User.Login,
		Title:     issue.Title,
		UpdatedAt: issue.UpdatedAt,
	}
}

func ConvertPR(pr *github_sdk.PullRequest, chatId int) ChangingDTO {
	return ChangingDTO{
		ChatId:    chatId,
		Link:      pr.URL,
		Event:     storage.PullRequest,
		Author:    pr.User.Login,
		Title:     pr.Title,
		UpdatedAt: pr.UpdatedAt,
	}
}
