package storage

type ChatStore interface {
	AddNewChat(chat *Chat) (int64, error)
	RemoveChat(id int64) error
	GetChatByID(id int64) (*Chat, error)
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
	RemoveRecord(chat_id int64, repo_id int) (int64, error)
	GetRecordByChat(chat_id int64) ([]ChatRepoRecord, error)
	GetRecordByLink(link_id int) ([]ChatRepoRecord, error)
	GetRecordById(id int) (*ChatRepoRecord, error)
	GetRecordOffset(start int, limit int) ([]ChatRepoRecord, error)
	GetRecordNumber() (int, error)
	GetRecordByChatAndLink(chatId int64, linkId int) (*ChatRepoRecord, error)
	UpdateRecord(record *ChatRepoRecord) error
}
