package server

import (
	"Common"
	"Crypto_Bot/MainServer/server/dtos"
	"Crypto_Bot/MainServer/storage"
	"time"
)

type StoreManager struct {
	chatStore           storage.ChatStore
	repoStore           storage.RepoStore
	chatRepoRecordStore storage.ChatRepoRecordStore
}

func NewStoreManager(chatStore storage.ChatStore, repoStore storage.RepoStore, chatRepoRecordStore storage.ChatRepoRecordStore) *StoreManager {
	return &StoreManager{chatStore, repoStore, chatRepoRecordStore}
}

func (sm *StoreManager) GetChats() ([]storage.Chat, error) {
	chatNumber, err := sm.chatStore.GetChatNumber()
	if err != nil {
		return nil, err
	}
	return sm.chatStore.GetChatsOffset(0, chatNumber)
}

func (sm *StoreManager) AddChat(chat *storage.Chat) (int64, error) {
	return sm.chatStore.AddNewChat(chat)
}

func (sm *StoreManager) DeleteChat(chatId int64) error {
	return sm.chatStore.RemoveChat(chatId)
}

func (sm *StoreManager) GetReposByChat(chatId int64) ([]storage.ChatRepoRecord, error) {
	return sm.chatRepoRecordStore.GetRecordByChat(chatId)
}

func (sm *StoreManager) AddRepo(repo *storage.Repo, repoDto *Common.RepoDTO, chatId int64) (int, error) {
	var id int
	needToUpdate := false
	oldRepo, err := sm.repoStore.GetRepoByOwnerAndName(repo.Owner, repo.Name)
	if oldRepo.Link == "" {
		id, err = sm.repoStore.AddNewRepo(repo)
		if err != nil {
			return -1, err
		}
	} else {
		if oldRepo.LastPR.IsZero() && !repo.LastPR.IsZero() {
			oldRepo.LastPR = repo.LastPR
			needToUpdate = true
		}
		if oldRepo.LastCommit.IsZero() && !repo.LastCommit.IsZero() {
			oldRepo.LastCommit = repo.LastCommit
			needToUpdate = true
		}
		if oldRepo.LastIssue.IsZero() && !repo.LastIssue.IsZero() {
			oldRepo.LastIssue = repo.LastIssue
			needToUpdate = true
		}
		if needToUpdate {
			err = sm.repoStore.UpdateRepo(oldRepo)
			if err != nil {
				return -1, err
			}
		}
		id = oldRepo.ID
	}
	chat, err := sm.chatStore.GetChatByID(chatId)
	if err != nil {
		return -1, nil
	}
	record, err := dtos.ParseChatRepoRecord(repoDto, chat, oldRepo)
	if err != nil {
		return -1, err
	}
	_, err = sm.chatRepoRecordStore.AddNewRecord(record)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (sm *StoreManager) DeleteRepo(chatId int64, owner string, name string) (int64, error) {
	repo, err := sm.repoStore.GetRepoByOwnerAndName(owner, name)
	if err != nil {
		return -1, err
	}
	num, err := sm.chatRepoRecordStore.RemoveRecord(chatId, repo.ID)
	if err != nil {
		return -1, err
	}
	records, err := sm.chatRepoRecordStore.GetRecordByLink(repo.ID)
	if err != nil {
		return -1, err
	}
	if len(records) == 0 {
		return 0, nil
	}
	checkPr := false
	checkCommit := false
	checkIssue := false
	for _, record := range records {
		if !record.Repo.LastCommit.IsZero() {
			checkCommit = true
		}
		if !record.Repo.LastIssue.IsZero() {
			checkIssue = true
		}
		if !record.Repo.LastPR.IsZero() {
			checkPr = true
		}
	}
	needUpdate := false
	if !checkPr {
		repo.LastPR = time.Time{}
		needUpdate = true
	}
	if !checkCommit {
		repo.LastCommit = time.Time{}
		needUpdate = true
	}
	if !checkIssue {
		repo.LastIssue = time.Time{}
		needUpdate = true
	}
	if needUpdate {
		err = sm.repoStore.UpdateRepo(repo)
		err = sm.repoStore.RemoveRepo(repo.ID)
		if err != nil {
			return -1, nil
		}
	}
	return num, nil
}

func (sm *StoreManager) GetReposByTag(chatId int64, tag string) ([]*storage.ChatRepoRecord, error) {
	records, err := sm.chatRepoRecordStore.GetRecordByChat(chatId)
	if err != nil {
		return nil, err
	}
	filter := make([]*storage.ChatRepoRecord, len(records))
	pointer := 0
	for _, record := range records {
		for _, tg := range record.Tags {
			if tg == tag {
				filter[pointer] = &record
				pointer++
				break
			}
		}
	}
	return filter[0:pointer], nil
}
