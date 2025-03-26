package storage

import "context"

type ChatStore interface {
	AddNewChat(ctx context.Context, chat *Chat) (int, error)
	RemoveChat(ctx context.Context, chat *Chat) error
	GetChatByID(ctx context.Context, id int) (*Chat, error)
	GetChatsOffset(ctx context.Context, start int, limit int) ([]Chat, error)
	GetChatNumber(ctx context.Context) (int, error)
	UpdateChat(ctx context.Context, chat *Chat) error
}

type RepoStore interface {
	AddNewRepo(ctx context.Context, repo *Repo) (int, error)
	RemoveRepo(ctx context.Context, repo *Repo) error
	GetRepoByID(ctx context.Context, id int) (*Repo, error)
	GetRepoByOwnerAndName(ctx context.Context, owner string, name string) (*Repo, error)
	GetReposOffset(ctx context.Context, start int, limit int) ([]Repo, error)
	GetRepoNumber(ctx context.Context) (int, error)
	UpdateRepo(ctx context.Context, repo *Repo) error
}

type ChatRepoRecordStore interface {
	AddNewRecord(ctx context.Context, record *ChatRepoRecord) (int, error)
	RemoveRecord(ctx context.Context, record *ChatRepoRecord) error
	GetRecordByChat(ctx context.Context, chat Chat) ([]ChatRepoRecord, error)
	GetRecordById(ctx context.Context, id int) (*ChatRepoRecord, error)
	GetRecordOffset(ctx context.Context, start int, limit int) ([]ChatRepoRecord, error)
	GetRecordNumber(ctx context.Context) (int, error)
	UpdateRecord(ctx context.Context, record *ChatRepoRecord) error
}
