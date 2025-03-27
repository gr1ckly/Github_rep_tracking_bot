package server

import "Crypto_Bot/MainServer/storage"

func ConvertRepoDTO(record *storage.ChatRepoRecord) (*RepoDTO, error) {
	ans := &RepoDTO{}
	ans.Link = record.Repo.Link
	ans.ChatID = record.Chat.ChatID
	ans.Events = record.Events
	ans.Tags = record.Tags
	return ans, nil
}

func ConvertChatDTO(chat *storage.Chat) (*ChatDTO, error) {
	ans := &ChatDTO{}
	ans.ChatID = chat.ChatID
	ans.Type = chat.Type
	return ans, nil
}
