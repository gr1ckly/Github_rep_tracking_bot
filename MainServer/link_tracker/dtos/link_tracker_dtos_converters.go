package dtos

import (
	"Common"
	"Crypto_Bot/MainServer/github_sdk"
	"strings"
)

func ConvertCommit(commit *github_sdk.Commit, chatId int64) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId: chatId,
		Link:   commit.HTMLURL[:strings.Index(commit.HTMLURL, "/commit/")],
		Event:  Common.Commit,
		Author: commit.Commit.Author.Name,
		Title:  commit.Commit.Message,
	}
}

func ConvertIssue(issue *github_sdk.Issue, chatId int64) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    chatId,
		Link:      issue.RepositoryURL,
		Event:     Common.Issue,
		Author:    issue.User.Login,
		Title:     issue.Title,
		UpdatedAt: issue.UpdatedAt,
	}
}

func ConvertPR(pr *github_sdk.PullRequest, chatId int64) Common.ChangingDTO {
	return Common.ChangingDTO{
		ChatId:    chatId,
		Link:      pr.Head.Repo.URL,
		Event:     Common.PullRequest,
		Author:    pr.User.Login,
		Title:     pr.Title,
		UpdatedAt: pr.UpdatedAt,
	}
}
