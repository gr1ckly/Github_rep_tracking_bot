package github_sdk

import "time"

type Commit struct {
	SHA         string     `json:"sha"`
	NodeID      string     `json:"node_id"`
	Commit      CommitInfo `json:"commit"`
	URL         string     `json:"url"`
	HTMLURL     string     `json:"html_url"`
	CommentsURL string     `json:"comments_url"`
	Author      User       `json:"author"`
	Committer   User       `json:"committer"`
	Parents     []Parent   `json:"parents"`
}

type Issue struct {
	ID                int             `json:"id"`
	NodeID            string          `json:"node_id"`
	URL               string          `json:"url"`
	RepositoryURL     string          `json:"repository_url"`
	LabelsURL         string          `json:"labels_url"`
	CommentsURL       string          `json:"comments_url"`
	EventsURL         string          `json:"events_url"`
	HTMLURL           string          `json:"html_url"`
	Number            int             `json:"number"`
	State             string          `json:"state"`
	Title             string          `json:"title"`
	Body              string          `json:"body"`
	User              User            `json:"user"`
	Labels            []Label         `json:"labels"`
	Assignee          User            `json:"assignee"`
	Assignees         []User          `json:"assignees"`
	Milestone         Milestone       `json:"milestone"`
	Locked            bool            `json:"locked"`
	ActiveLockReason  string          `json:"active_lock_reason"`
	Comments          int             `json:"comments"`
	PullRequest       PullRequestInfo `json:"pull_request"`
	ClosedAt          *time.Time      `json:"closed_at"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	ClosedBy          User            `json:"closed_by"`
	AuthorAssociation string          `json:"author_association"`
	StateReason       string          `json:"state_reason"`
}

type PullRequest struct {
	ID                  int        `json:"id"`
	NodeID              string     `json:"node_id"`
	URL                 string     `json:"url"`
	HTMLURL             string     `json:"html_url"`
	DiffURL             string     `json:"diff_url"`
	PatchURL            string     `json:"patch_url"`
	IssueURL            string     `json:"issue_url"`
	CommitsURL          string     `json:"commits_url"`
	ReviewCommentsURL   string     `json:"review_comments_url"`
	ReviewCommentURL    string     `json:"review_comment_url"`
	CommentsURL         string     `json:"comments_url"`
	StatusesURL         string     `json:"statuses_url"`
	Number              int        `json:"number"`
	State               string     `json:"state"`
	Title               string     `json:"title"`
	Body                string     `json:"body"`
	User                User       `json:"user"`
	Labels              []Label    `json:"labels"`
	Assignee            *User      `json:"assignee"`
	Assignees           []User     `json:"assignees"`
	Milestone           *Milestone `json:"milestone"`
	Locked              bool       `json:"locked"`
	ActiveLockReason    string     `json:"active_lock_reason"`
	CreatedAt           string     `json:"created_at"`
	UpdatedAt           string     `json:"updated_at"`
	ClosedAt            *string    `json:"closed_at"`
	MergedAt            *string    `json:"merged_at"`
	MergeCommitSHA      *string    `json:"merge_commit_sha"`
	AuthorAssociation   string     `json:"author_association"`
	Draft               bool       `json:"draft"`
	Head                Branch     `json:"head"`
	Base                Branch     `json:"base"`
	Merged              bool       `json:"merged"`
	Mergeable           *bool      `json:"mergeable"`
	Rebaseable          *bool      `json:"rebaseable"`
	MergeableState      string     `json:"mergeable_state"`
	Comments            int        `json:"comments"`
	ReviewComments      int        `json:"review_comments"`
	MaintainerCanModify bool       `json:"maintainer_can_modify"`
	Commits             int        `json:"commits"`
	Additions           int        `json:"additions"`
	Deletions           int        `json:"deletions"`
	ChangedFiles        int        `json:"changed_files"`
}

type CommitInfo struct {
	Author       CommitAuthor `json:"author"`
	Committer    CommitAuthor `json:"committer"`
	Message      string       `json:"message"`
	Tree         Tree         `json:"tree"`
	URL          string       `json:"url"`
	CommentCount int          `json:"comment_count"`
	Verification Verification `json:"verification"`
}

type CommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"url"`
	HTMLURL   string `json:"html_url"`
}

type Tree struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type Parent struct {
	SHA     string `json:"sha"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}

type Verification struct {
	Verified  bool   `json:"verified"`
	Reason    string `json:"reason"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
}

type Label struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

type Milestone struct {
	ID           int        `json:"id"`
	NodeID       string     `json:"node_id"`
	Number       int        `json:"number"`
	State        string     `json:"state"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Creator      User       `json:"creator"`
	OpenIssues   int        `json:"open_issues"`
	ClosedIssues int        `json:"closed_issues"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ClosedAt     *time.Time `json:"closed_at"`
	DueOn        *time.Time `json:"due_on"`
}

type PullRequestInfo struct {
	URL      string `json:"url"`
	HTMLURL  string `json:"html_url"`
	DiffURL  string `json:"diff_url"`
	PatchURL string `json:"patch_url"`
}

type Branch struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	SHA   string `json:"sha"`
	User  User   `json:"user"`
	Repo  Repo   `json:"repo"`
}

type Repo struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    User   `json:"owner"`
	Private  bool   `json:"private"`
	URL      string `json:"url"`
}
