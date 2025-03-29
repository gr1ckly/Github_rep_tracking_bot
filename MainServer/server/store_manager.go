package server

import (
	"Crypto_Bot/MainServer/storage"
	"context"
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

func (sm *StoreManager) GetChats(ctx context.Context) ([]storage.Chat, error) {
	chatNumber, err := sm.chatStore.GetChatNumber(ctx)
	if err != nil {
		return nil, err
	}
	return sm.chatStore.GetChatsOffset(ctx, 0, chatNumber)
}

func (sm *StoreManager) AddChat(ctx context.Context, chat *storage.Chat) (int, error) {
	return sm.chatStore.AddNewChat(ctx, chat)
}

func (sm *StoreManager) DeleteChat(ctx context.Context, chatId int) error {
	return sm.chatStore.RemoveChat(ctx, chatId)
}

func (sm *StoreManager) GetReposByChat(ctx context.Context, chatId int) ([]storage.ChatRepoRecord, error) {
	return sm.chatRepoRecordStore.GetRecordByChat(ctx, chatId)
}

func (sm *StoreManager) AddRepo(ctx context.Context, repo *storage.Repo, repoDto *RepoDTO) (int, error) {
	var id int
	needToUpdate := false
	oldRepo, err := sm.repoStore.GetRepoByOwnerAndName(ctx, repo.Owner, repo.Name)
	if oldRepo.Link == "" {
		id, err = sm.repoStore.AddNewRepo(ctx, repo)
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
			err = sm.repoStore.UpdateRepo(ctx, oldRepo)
			if err != nil {
				return -1, err
			}
		}
		id = oldRepo.ID
	}
	chat, err := sm.chatStore.GetChatByID(ctx, repoDto.ChatID)
	if err != nil {
		return -1, nil
	}
	record, err := ParseChatRepoRecord(repoDto, chat, oldRepo)
	if err != nil {
		return -1, err
	}
	_, err = sm.chatRepoRecordStore.AddNewRecord(ctx, record)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (sm *StoreManager) DeleteRepo(ctx context.Context, chatId int, owner string, name string) error {
	repo, err := sm.repoStore.GetRepoByOwnerAndName(ctx, owner, name)
	if err != nil {
		return err
	}
	err = sm.chatRepoRecordStore.RemoveRecord(ctx, chatId, repo.ID)
	if err != nil {
		return err
	}
	records, err := sm.chatRepoRecordStore.GetRecordByLink(ctx, repo.ID)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		err = sm.repoStore.RemoveRepo(ctx, repo.ID)
		if err != nil {
			return err
		}
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
		err = sm.repoStore.UpdateRepo(ctx, repo)
		err = sm.repoStore.RemoveRepo(ctx, repo.ID)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (sm *StoreManager) GetReposByTag(ctx context.Context, chatId int, tag string) ([]*storage.ChatRepoRecord, error) {
	records, err := sm.chatRepoRecordStore.GetRecordByChat(ctx, chatId)
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
