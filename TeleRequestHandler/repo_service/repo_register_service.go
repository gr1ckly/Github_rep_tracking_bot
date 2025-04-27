package repo_service

import (
	"Common"
)

type RepoRegisterService interface {
	AddRepo(chatId int64, dto Common.RepoDTO) error
	GetReposByChat(chatId int64) ([]Common.RepoDTO, error)
	GetReposByTag(chatId int64, tag string) ([]Common.RepoDTO, error)
	DeleteRepo(chatId int64, link string) error
}
