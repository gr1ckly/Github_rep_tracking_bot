package github_sdk

type GithubService interface {
	GetCommits(string, string) ([]Commit, error)
	GetIssues(string, string) ([]Issue, error)
	GetPullRequests(string, string) ([]PullRequest, error)
}
