package state_machine

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateName int

type State struct {
	Name StateName
	StateHandler
}

func NewState(stateName StateName, handler StateHandler) State {
	return State{stateName, handler}
}

type StateHandler interface {
	Start(usrCtx UserContext) (UserContext, error)
	Process(usrCtx UserContext, update tgbotapi.Update) (UserContext, error)
}

const (
	NONE StateName = iota
	WAIT_LINK
	WAIT_EVENTS
	WAIT_TAGS
)
