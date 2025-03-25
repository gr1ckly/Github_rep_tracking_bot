package github_sdk

import (
	"Crypto_Bot/Scrapper/custom_errors"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

type HttpGithubService struct {
	apiUrlMutex  sync.RWMutex
	apiUrl       string
	token        string
	acceptFormat string
	apiVersion   string
	client       http.Client
}

func NewHttpGithubService(apiUrl string, token string, acceptFormat string, apiVersion string, timeout int) *HttpGithubService {
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	return &HttpGithubService{apiUrl: apiUrl, token: token, acceptFormat: acceptFormat, apiVersion: apiVersion, client: client}
}

func (ghService *HttpGithubService) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", ghService.token)
	req.Header.Set("Accept", ghService.acceptFormat)
	req.Header.Set("X-GitHub-Api-Version", ghService.apiVersion)
}

func (ghService *HttpGithubService) fetch(repoName string, owner string, way string) ([]byte, error) {
	ghService.apiUrlMutex.RLock()
	req, err := http.NewRequest("GET", ghService.apiUrl+"/repos/"+owner+"/"+repoName+"/"+way, nil)
	ghService.apiUrlMutex.RUnlock()
	if err != nil {
		return nil, err
	}
	ghService.setHeaders(req)
	resp, err := ghService.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 301 {
		ghService.apiUrl = resp.Header.Get("Location")
		return ghService.fetch(repoName, owner, way)
	}
	if resp.StatusCode == 302 || resp.StatusCode == 307 {
		prev := ghService.apiUrl
		ghService.apiUrlMutex.Lock()
		ghService.apiUrl = resp.Header.Get("Location")
		ghService.apiUrlMutex.Unlock()
		ans, err := ghService.fetch(repoName, owner, way)
		ghService.apiUrlMutex.Lock()
		ghService.apiUrl = prev
		ghService.apiUrlMutex.Unlock()
		return ans, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, custom_errors.StatusError{resp.StatusCode, resp.Status}
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (ghService *HttpGithubService) GetCommits(repoName string, owner string) ([]Commit, error) {
	var commits []Commit
	data, err := ghService.fetch(repoName, owner, "commits")
	err = json.Unmarshal(data, commits)
	return commits, err
}
func (ghService *HttpGithubService) GetIssues(repoName string, owner string) ([]Issue, error) {
	var issues []Issue
	data, err := ghService.fetch(repoName, owner, "commits")
	err = json.Unmarshal(data, issues)
	return issues, err
}
func (ghService *HttpGithubService) GetPullRequests(repoName string, owner string) ([]PullRequest, error) {
	var pullRequests []PullRequest
	data, err := ghService.fetch(repoName, owner, "commits")
	err = json.Unmarshal(data, pullRequests)
	return pullRequests, err
}
