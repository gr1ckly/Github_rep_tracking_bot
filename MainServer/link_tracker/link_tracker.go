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
	"slices"
	"sync"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

type LinkTracker struct {
	ghServiceMutex    sync.Mutex
	ghService         github_sdk.GithubService
	storeManagerMutex sync.Mutex
	storeManager      *server.StoreManager
	recordStoreMutex  sync.Mutex
	recordsStore      storage.ChatRepoRecordStore
	repoStoreMutex    sync.Mutex
	repoStore         storage.RepoStore
	batchSize         int
	observers         []Observer[any]
	scheduler         gocron.Scheduler
}

func NewLinkTracker(ghService github_sdk.GithubService, storeManager *server.StoreManager, recordsStore storage.ChatRepoRecordStore, repoStore storage.RepoStore, batchSize int, observers ...Observer[any]) (*LinkTracker, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &LinkTracker{sync.Mutex{}, ghService, sync.Mutex{}, storeManager, sync.Mutex{}, recordsStore, sync.Mutex{}, repoStore, batchSize, observers, scheduler}, nil
}

func (lt LinkTracker) checkAllLinks() interface{} {
	lt.repoStoreMutex.Lock()
	repoNumber, err := lt.repoStore.GetRepoNumber()
	lt.repoStoreMutex.Unlock()
	if err != nil {
		return nil
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()
	for i := 0; i < repoNumber; i += lt.batchSize {
		go func(start int, limit int) {
			wg.Add(1)
			defer wg.Done()
			lt.repoStoreMutex.Lock()
			repos, _ := lt.repoStore.GetReposOffset(start, limit)
			lt.repoStoreMutex.Unlock()
			for _, repo := range repos {
				err = lt.checkLink(&repo)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}(i, repoNumber)
	}
	return nil
}

func (lt LinkTracker) tryCheckStatusError(repo *storage.Repo, err error) error {
	var statusErr Common.StatusError
	if errors.As(err, &statusErr) {
		if statusErr.StatusCode == 404 {
			lt.repoStoreMutex.Lock()
			defer lt.repoStoreMutex.Unlock()
			err = lt.repoStore.RemoveRepo(repo.ID)
			if err != nil {
				return statusErr
			}
		}
	}
	return err
}

func (lt LinkTracker) checkLink(repo *storage.Repo) error {
	needToUpdateRepo := false
	lt.recordStoreMutex.Lock()
	records, err := lt.recordsStore.GetRecordByLink(repo.ID)
	defer lt.recordStoreMutex.Unlock()
	if err != nil {
		return err
	}
	if !repo.LastCommit.IsZero() {
		lt.ghServiceMutex.Lock()
		commits, err := lt.ghService.GetCommits(repo.Name, repo.Owner, repo.LastCommit)
		lt.ghServiceMutex.Unlock()
		if err != nil {
			return lt.tryCheckStatusError(repo, err)
		}
		for _, commit := range commits {
			if commit.Commit.Committer.Date.After(repo.LastCommit) {
				repo.LastCommit = commit.Commit.Committer.Date.Add(1 * time.Second)
				needToUpdateRepo = true
			}
			for _, record := range records {
				if slices.Contains(record.Events, Common.Commit) {
					newCommitChange := dtos.ConvertCommit(&commit, record.Chat.ChatID, repo.Link)
					lt.NotifyAll(newCommitChange)
				}
			}
		}
	}
	if !repo.LastIssue.IsZero() {
		lt.ghServiceMutex.Lock()
		issues, err := lt.ghService.GetIssues(repo.Name, repo.Owner, repo.LastIssue)
		lt.ghServiceMutex.Unlock()
		if err != nil {
			return lt.tryCheckStatusError(repo, err)
		}
		for _, issue := range issues {
			for _, record := range records {
				if slices.Contains(record.Events, Common.Issue) {
					newCommitChange := dtos.ConvertIssue(&issue, record.Chat.ChatID, repo.Link)
					lt.NotifyAll(newCommitChange)
				}
			}
			if issue.UpdatedAt.After(repo.LastCommit) {
				repo.LastIssue = issue.UpdatedAt.Add(1 * time.Second)
				needToUpdateRepo = true
			}
		}
	}
	if !repo.LastPR.IsZero() {
		lt.ghServiceMutex.Lock()
		prs, err := lt.ghService.GetPullRequests(repo.Name, repo.Owner, repo.LastPR)
		lt.ghServiceMutex.Unlock()
		if err != nil {
			return lt.tryCheckStatusError(repo, err)
		}
		for _, pr := range prs {
			if pr.UpdatedAt.After(repo.LastCommit) {
				repo.LastPR = pr.UpdatedAt.Add(1 * time.Second)
				needToUpdateRepo = true
			}
			for _, record := range records {
				if slices.Contains(record.Events, Common.PullRequest) {
					newCommitChange := dtos.ConvertPR(&pr, record.Chat.ChatID, repo.Link)
					lt.NotifyAll(newCommitChange)
				}
			}
		}
	}
	if needToUpdateRepo {
		lt.repoStoreMutex.Lock()
		defer lt.repoStoreMutex.Unlock()
		return lt.repoStore.UpdateRepo(repo)
	}
	return nil
}

func (lt LinkTracker) StartTracking() {
	lt.scheduler.NewJob(gocron.CronJob("* * * * *", true), gocron.NewTask(lt.checkAllLinks), gocron.WithSingletonMode(gocron.LimitModeReschedule))
	lt.scheduler.Start()
}

func (lt LinkTracker) Stop() error {
	return lt.scheduler.Shutdown()
}

func (lt LinkTracker) NotifyAll(msg any) {
	for _, obs := range lt.observers {
		err := obs.Notify(msg)
		if err != nil {
			logger.Warn("Writing message err: " + err.Error())
		}
	}
}

func (lt LinkTracker) AddObserver(observer Observer[any]) {
	lt.observers = append(lt.observers, observer)
}

func (lt LinkTracker) RemoveObserver(observer Observer[any]) {
	newObservers := []Observer[any]{}
	for _, obs := range lt.observers {
		if obs != observer {
			newObservers = append(newObservers, obs)
		}
	}
	lt.observers = newObservers
}
