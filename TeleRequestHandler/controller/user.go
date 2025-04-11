package controller

import "TeleRequestHandler/controller/state_machine"

type User struct {
	ChatId int
	State  state_machine.State
}

func NewUser(chatId int) User {
	return User{chatId, state_machine.UNAUTHORIZED}
}
