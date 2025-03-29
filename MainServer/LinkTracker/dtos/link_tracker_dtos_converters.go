package dtos

import (
	"Crypto_Bot/MainServer/github_sdk"
)

func ConvertCommit(commit *github_sdk.Commit, chatId int) CommitChangesDTO {
	return CommitChangesDTO{
		ChatId:   chatId,
		Link:     commit.URL,
		Author:   commit.Commit.Author.Name,
		Commiter: commit.Commit.Committer.Name,
		Message:  commit.Commit.Message,
		Branch:   commit.Commit.Tree.URL,
	}
}

func ConvertIssue(issue *github_sdk.Issue, chatId int) IssueChangesSTO {
	return IssueChangesSTO{
		ChatId:    chatId,
		Link:      issue.URL,
		Author:    issue.User.Login,
		Title:     issue.Title,
		UpdatedAt: issue.UpdatedAt,
	}
}

func ConvertPR(pr *github_sdk.PullRequest, chatId int) PullRequestChangesDTO {
	return PullRequestChangesDTO{
		ChatId:    chatId,
		Link:      pr.URL,
		Author:    pr.User.Login,
		Number:    pr.Number,
		Title:     pr.Title,
		Status:    pr.State,
		UpdatedAt: pr.UpdatedAt,
	}
}
