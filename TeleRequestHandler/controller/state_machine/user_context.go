package state_machine

type UserContext struct {
	ChatId       int64
	CurrentState *State
	CommandName  string
	Tags         []string
	Link         string
	Events       []string
}

func GetDefaultContext(chatId int64) *UserContext {
	return &UserContext{ChatId: chatId}
}
