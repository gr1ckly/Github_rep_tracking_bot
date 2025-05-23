package Common

type ChatDTO struct {
	ChatID int64  `json:"chat_id"`
	Type   string `json:"type"`
}

type RepoDTO struct {
	Link   string   `json:"link"`
	Tags   []string `json:"tags"`
	Events []string `json:"events"`
}

type ErrorDTO struct {
	Error string `json:"error"`
}

type ResultDTO[T any] struct {
	Result any `json:"result"`
}
