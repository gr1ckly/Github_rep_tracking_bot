package storage

type ChatStore interface {
	AddNewChat(chat *Chat) (int, error)
	RemoveChat(id int) error
	GetChatByID(id int) (*Chat, error)
	GetChatsOffset(start int, limit int) ([]Chat, error)
	GetChatNumber() (int, error)
	UpdateChat(chat *Chat) error
}

type RepoStore interface {
	AddNewRepo(repo *Repo) (int, error)
	RemoveRepo(id int) error
	GetRepoByID(id int) (*Repo, error)
	GetRepoByOwnerAndName(owner string, name string) (*Repo, error)
	GetReposOffset(start int, limit int) ([]Repo, error)
	GetRepoNumber() (int, error)
	UpdateRepo(repo *Repo) error
}

type ChatRepoRecordStore interface {
	AddNewRecord(record *ChatRepoRecord) (int, error)
	RemoveRecord(chat_id int, repo_id int) error
	GetRecordByChat(chat_id int) ([]ChatRepoRecord, error)
	GetRecordByLink(link_id int) ([]ChatRepoRecord, error)
	GetRecordById(id int) (*ChatRepoRecord, error)
	GetRecordOffset(start int, limit int) ([]ChatRepoRecord, error)
	GetRecordNumber() (int, error)
	UpdateRecord(record *ChatRepoRecord) error
}
