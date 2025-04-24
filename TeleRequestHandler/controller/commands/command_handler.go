package commands

import (
	"TeleRequestHandler/controller/state_machine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler interface {
	Execute(usrCtx state_machine.UserContext, upd tgbotapi.Update) state_machine.UserContext
}
