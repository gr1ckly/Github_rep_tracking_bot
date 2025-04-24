package link_tracker

import (
	"Common"
	"Crypto_Bot/MainServer/github_sdk"
	"Crypto_Bot/MainServer/link_tracker/dtos"
	"Crypto_Bot/MainServer/server"
	"Crypto_Bot/MainServer/storage"
	"errors"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"os"
	"sync"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

type LinkTracker struct {
	ghService    github_sdk.GithubService
	storeManager *server.StoreManager
	recordsStore storage.ChatRepoRecordStore
	repoStore    storage.RepoStore
	batchSize    int
	observers    []Observer[any]
	scheduler    gocron.Scheduler
}

func NewLinkTracker(ghService github_sdk.GithubService, storeManager *server.StoreManager, recordsStore storage.ChatRepoRecordStore, repoStore storage.RepoStore, batchSize int, observers ...Observer[any]) (*LinkTracker, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &LinkTracker{ghService, storeManager, recordsStore, repoStore, batchSize, observers, scheduler}, nil
}

func (lt *LinkTracker) checkAllLinks() interface{} {
	chatNumber, err := lt.recordsStore.GetRecordNumber()
	if err != nil {
		return nil
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()
	for i := 0; i < chatNumber; i += lt.batchSize {
		go func(start int, limit int) {
			wg.Add(1)
			defer wg.Done()
			records, _ := lt.recordsStore.GetRecordOffset(start, limit)
			for _, record := range records {
				err = lt.checkLink(&record)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}(i, chatNumber)
	}
	return nil
}

func (lt *LinkTracker) tryCheckStatusError(record *storage.ChatRepoRecord, err error) error {
	var statusErr Common.StatusError
	if errors.As(err, &statusErr) {
		if statusErr.StatusCode == 404 {
			_, err = lt.storeManager.DeleteRepo(record.Chat.ChatID, record.Repo.Owner, record.Repo.Name)
			if err != nil {
				return statusErr
			}
		}
	}
	return err
}

func (lt *LinkTracker) checkLink(record *storage.ChatRepoRecord) error {
	needToUpdateRepo := false
	if !record.Repo.LastCommit.IsZero() {
		commits, err := lt.ghService.GetCommits(record.Repo.Name, record.Repo.Owner, record.Repo.LastCommit)
		if err != nil {
			return lt.tryCheckStatusError(record, err)
		}
		for _, commit := range commits {
			newCommitChange := dtos.ConvertCommit(&commit, record)
			if commit.Committer.Date.After(record.Repo.LastCommit) {
				record.Repo.LastCommit = commit.Committer.Date.Add(1 * time.Second)
				needToUpdateRepo = true
			}
			lt.NotifyAll(newCommitChange)
		}
	}
	if !record.Repo.LastIssue.IsZero() {
		issues, err := lt.ghService.GetIssues(record.Repo.Name, record.Repo.Owner, record.Repo.LastIssue)
		if err != nil {
			return lt.tryCheckStatusError(record, err)
		}
		for _, issue := range issues {
			newIssueChange := dtos.ConvertIssue(&issue, record)
			if issue.UpdatedAt.After(record.Repo.LastCommit) {
				record.Repo.LastIssue = issue.UpdatedAt.Add(1 * time.Second)
				needToUpdateRepo = true
			}
			lt.NotifyAll(newIssueChange)
		}
	}
	if !record.Repo.LastPR.IsZero() {
		prs, err := lt.ghService.GetPullRequests(record.Repo.Name, record.Repo.Owner, record.Repo.LastPR)
		if err != nil {
			return lt.tryCheckStatusError(record, err)
		}
		for _, pr := range prs {
			newPRChange := dtos.ConvertPR(&pr, record)
			if pr.UpdatedAt.After(record.Repo.LastCommit) {
				record.Repo.LastPR = pr.UpdatedAt.Add(1 * time.Second)
				needToUpdateRepo = true
			}
			lt.NotifyAll(newPRChange)
		}
	}
	if needToUpdateRepo {
		return lt.repoStore.UpdateRepo(record.Repo)
	}
	return nil
}

func (lt *LinkTracker) StartTracking() {
	lt.scheduler.NewJob(gocron.CronJob("* * * * *", true), gocron.NewTask(lt.checkAllLinks()), gocron.WithSingletonMode(gocron.LimitModeReschedule))
	lt.scheduler.Start()
}

func (lt *LinkTracker) Stop() error {
	return lt.scheduler.Shutdown()
}

func (lt *LinkTracker) NotifyAll(msg any) {
	for _, obs := range lt.observers {
		err := obs.Notify(msg)
		if err != nil {
			logger.Warn("Writing message err: " + err.Error())
		}
	}
}

func (lt *LinkTracker) AddObserver(observer Observer[any]) {
	lt.observers = append(lt.observers, observer)
}

func (lt *LinkTracker) RemoveObserver(observer Observer[any]) {
	newObservers := []Observer[any]{}
	for _, obs := range lt.observers {
		if obs != observer {
			newObservers = append(newObservers, obs)
		}
	}
	lt.observers = newObservers
}
