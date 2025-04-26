package dtos

import (
	"Common"
	"Crypto_Bot/MainServer/storage"
)

func ConvertRepoDTO(record *storage.ChatRepoRecord) (*Common.RepoDTO, error) {
	ans := &Common.RepoDTO{}
	ans.Link = record.Repo.Link
	ans.Events = record.Events
	ans.Tags = record.Tags
	return ans, nil
}

func ConvertChatDTO(chat *storage.Chat) (*Common.ChatDTO, error) {
	ans := &Common.ChatDTO{}
	ans.ChatID = chat.ChatID
	ans.Type = chat.Type
	return ans, nil
}
