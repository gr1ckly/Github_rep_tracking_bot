package repo_service

import (
	"Common"
)

type RepoRegisterService interface {
	AddRepo(chatId int, dto Common.RepoDTO) error
	GetReposByChat(chatId int) ([]Common.RepoDTO, error)
	GetReposByTag(chatId int, tag string) ([]Common.RepoDTO, error)
	DeleteRepo(chatId int, link string) error
}
