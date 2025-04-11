package dtos

import (
	"Common"
	"Crypto_Bot/MainServer/custom_errors"
	"Crypto_Bot/MainServer/storage"
	"strings"
	"time"
)

func ParseNameAndOwner(link string) (string, string, error) {
	ownerName := strings.TrimPrefix(link, "https://github.com/")
	splitLine := strings.Split(ownerName, "/")
	if len(splitLine) < 2 {
		return "", "", custom_errors.InvalidLinkError{link}
	}
	return splitLine[1], splitLine[0], nil
}

func ParseRepo(dto Common.RepoDTO) (*storage.Repo, error) {
	ans := &storage.Repo{}
	ans.Link = dto.Link
	name, owner, err := ParseNameAndOwner(ans.Link)
	if err != nil {
		return nil, err
	}
	ans.Name = name
	ans.Owner = owner
	if len(dto.Events) == 0 {
		ans.LastCommit = time.Now()
		ans.LastIssue = time.Now()
		ans.LastPR = time.Now()
	} else {
		for _, event := range dto.Events {
			switch event {
			case string(Common.PullRequest):
				ans.LastPR = time.Now()
			case string(Common.Commit):
				ans.LastCommit = time.Now()
			case string(Common.Issue):
				ans.LastIssue = time.Now()
			default:
				return nil, custom_errors.InvalidRepoEventsError{dto.Events}
			}
		}
	}
	return ans, nil
}

func ParseChat(dto Common.ChatDTO) (*storage.Chat, error) {
	ans := &storage.Chat{}
	ans.ChatID = dto.ChatID
	ans.Type = dto.Type
	return ans, nil
}

func ParseChatRepoRecord(repoDto *Common.RepoDTO, chat *storage.Chat, repo *storage.Repo) (*storage.ChatRepoRecord, error) {
	ans := &storage.ChatRepoRecord{Repo: repo, Chat: chat, Tags: repoDto.Tags, Events: repoDto.Events}
	return ans, nil
}
