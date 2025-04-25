package server

import (
	"Common"
	"Crypto_Bot/MainServer/custom_errors"
	"Crypto_Bot/MainServer/server/dtos"
	"Crypto_Bot/MainServer/storage"
	"fmt"
	"sync"
	"time"
)

type StoreManager struct {
	chatStoreMutex      sync.Mutex
	chatStore           storage.ChatStore
	repoStoreMutex      sync.Mutex
	repoStore           storage.RepoStore
	recordMutex         sync.Mutex
	chatRepoRecordStore storage.ChatRepoRecordStore
}

func NewStoreManager(chatStore storage.ChatStore, repoStore storage.RepoStore, chatRepoRecordStore storage.ChatRepoRecordStore) *StoreManager {
	return &StoreManager{sync.Mutex{}, chatStore, sync.Mutex{}, repoStore, sync.Mutex{}, chatRepoRecordStore}
}

func (sm *StoreManager) GetChats() ([]storage.Chat, error) {
	sm.chatStoreMutex.Lock()
	defer sm.chatStoreMutex.Unlock()
	chatNumber, err := sm.chatStore.GetChatNumber()
	if err != nil {
		return nil, err
	}
	ans, err := sm.chatStore.GetChatsOffset(0, chatNumber)
	return ans, err
}

func (sm *StoreManager) AddChat(chat *storage.Chat) (int64, error) {
	sm.chatStoreMutex.Lock()
	defer sm.chatStoreMutex.Unlock()
	id, err := sm.chatStore.AddNewChat(chat)
	return id, err
}

func (sm *StoreManager) DeleteChat(chatId int64) error {
	sm.chatStoreMutex.Lock()
	defer sm.chatStoreMutex.Unlock()
	return sm.chatStore.RemoveChat(chatId)
}

func (sm *StoreManager) GetReposByChat(chatId int64) ([]storage.ChatRepoRecord, error) {
	sm.recordMutex.Lock()
	defer sm.recordMutex.Unlock()
	ans, err := sm.chatRepoRecordStore.GetRecordByChat(chatId)
	return ans, err
}

func (sm *StoreManager) AddRepo(repo *storage.Repo, repoDto *Common.RepoDTO, chatId int64) (int, error) {
	var id int
	needToUpdate := false
	sm.repoStoreMutex.Lock()
	defer sm.repoStoreMutex.Unlock()
	oldRepo, err := sm.repoStore.GetRepoByOwnerAndName(repo.Owner, repo.Name)
	if oldRepo == nil {
		id, err = sm.repoStore.AddNewRepo(repo)
		if err != nil {
			return -1, err
		} else {
			return id, nil
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
	sm.chatStoreMutex.Lock()
	defer sm.chatStoreMutex.Unlock()
	chat, err := sm.chatStore.GetChatByID(chatId)
	if err != nil {
		return -1, nil
	}
	record, err := dtos.ParseChatRepoRecord(repoDto, chat, oldRepo)
	if err != nil {
		return -1, err
	}
	sm.recordMutex.Lock()
	defer sm.recordMutex.Unlock()
	oldRecord, err := sm.chatRepoRecordStore.GetRecordByChatAndLink(chatId, oldRepo.ID)
	if oldRecord != nil {
		return -1, custom_errors.NewAlreadyExistsError(fmt.Errorf("Repo already tracking"))
	}
	_, err = sm.chatRepoRecordStore.AddNewRecord(record)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (sm *StoreManager) DeleteRepo(chatId int64, owner string, name string) (int64, error) {
	sm.repoStoreMutex.Lock()
	defer sm.repoStoreMutex.Unlock()
	repo, err := sm.repoStore.GetRepoByOwnerAndName(owner, name)
	if err != nil {
		return -1, err
	}
	sm.recordMutex.Lock()
	defer sm.recordMutex.Unlock()
	num, err := sm.chatRepoRecordStore.RemoveRecord(chatId, repo.ID)
	if err != nil {
		return -1, err
	}
	records, err := sm.chatRepoRecordStore.GetRecordByLink(repo.ID)
	if err != nil {
		return -1, err
	}
	if records == nil || len(records) == 0 {
		return num, nil
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
	sm.recordMutex.Lock()
	defer sm.recordMutex.Unlock()
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
