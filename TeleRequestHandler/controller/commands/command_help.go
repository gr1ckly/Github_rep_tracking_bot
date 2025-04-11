package commands

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/state_machine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type CommandHelpHandler struct {
	bot bot.Bot[any, string, int64]
}

func NewCommandHelpHandler(bot bot.Bot[any, string, int64]) *CommandHelpHandler {
	return &CommandHelpHandler{bot}
}

func (cs *CommandHelpHandler) Execute(usrCtx *state_machine.UserContext, upd tgbotapi.Update) {
	if usrCtx.CurrentState.Name == state_machine.NONE {
		cmdMap := bot.GetCommandsDescription()
		builder := strings.Builder{}
		for key, _ := range cmdMap {
			builder.WriteString(key)
			builder.WriteString(" : ")
			builder.WriteString(cmdMap[key])
			builder.WriteRune('\n')
		}
		err := cs.bot.SendMessage(usrCtx.ChatId, builder.String())
		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		err := cs.bot.SendMessage(usrCtx.ChatId, "Недоступная команда")
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
}
