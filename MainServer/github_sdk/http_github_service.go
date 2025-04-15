package github_sdk

import (
	"Common"
	"Crypto_Bot/MainServer/server/dtos"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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

func (ghService *HttpGithubService) fetch(repoName string, owner string, way string, since *time.Time) ([]byte, error) {
	ghService.apiUrlMutex.RLock()
	baseUrl := ghService.apiUrl + "repos/" + owner + "/" + repoName
	if way != "" {
		baseUrl += "/" + way
	}
	req, err := http.NewRequest("GET", baseUrl, nil)
	ghService.apiUrlMutex.RUnlock()
	if err != nil {
		return nil, err
	}
	ghService.setHeaders(req)
	if since != nil {
		params := url.Values{}
		params.Add("since", since.Format(time.RFC3339))
		req.URL.RawQuery = params.Encode()
	}
	resp, err := ghService.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 301 {
		ghService.apiUrlMutex.Lock()
		ghService.apiUrl = resp.Header.Get("Location")
		ghService.apiUrlMutex.Unlock()
		return ghService.fetch(repoName, owner, way, since)
	}
	if resp.StatusCode == 302 || resp.StatusCode == 307 {
		prev := ghService.apiUrl
		ghService.apiUrlMutex.Lock()
		ghService.apiUrl = resp.Header.Get("Location")
		ghService.apiUrlMutex.Unlock()
		ans, err := ghService.fetch(repoName, owner, way, since)
		ghService.apiUrlMutex.Lock()
		ghService.apiUrl = prev
		ghService.apiUrlMutex.Unlock()
		return ans, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, Common.StatusError{resp.StatusCode, resp.Status, resp.Request.RequestURI}
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (ghService *HttpGithubService) GetCommits(repoName string, owner string, since time.Time) ([]Commit, error) {
	var commits []Commit
	data, err := ghService.fetch(repoName, owner, "commits", &since)
	if data == nil {
		return nil, err
	}
	err = json.Unmarshal(data, &commits)
	return nil, err
}
func (ghService *HttpGithubService) GetIssues(repoName string, owner string, since time.Time) ([]Issue, error) {
	var issues []Issue
	data, err := ghService.fetch(repoName, owner, "issues", &since)
	if data == nil {
		return nil, err
	}
	err = json.Unmarshal(data, &issues)
	return issues, err
}
func (ghService *HttpGithubService) GetPullRequests(repoName string, owner string, since time.Time) ([]PullRequest, error) {
	var pullRequests []PullRequest
	data, err := ghService.fetch(repoName, owner, "pulls", &since)
	if data == nil {
		return nil, err
	}
	err = json.Unmarshal(data, &pullRequests)
	return pullRequests, err
}

func (ghService *HttpGithubService) Check(link string) bool {
	name, owner, err := dtos.ParseNameAndOwner(link)
	if err != nil {
		return false
	}
	_, err = ghService.fetch(name, owner, "", nil)
	return err == nil
}
