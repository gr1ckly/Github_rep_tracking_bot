package commands

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/state_machine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type CommandHelpHandler struct {
	bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
}

func NewCommandHelpHandler(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) *CommandHelpHandler {
	return &CommandHelpHandler{bot}
}

func (cs *CommandHelpHandler) Execute(usrCtx *state_machine.UserContext, upd tgbotapi.Update) {
	if usrCtx.CurrentState.Name == state_machine.NONE {
		cmdMap := bot.GetCommandsDescription()
		builder := strings.Builder{}
		for key, _ := range cmdMap {
			builder.WriteRune('/')
			builder.WriteString(key)
			builder.WriteString(" : ")
			builder.WriteString(cmdMap[key])
			builder.WriteRune('\n')
		}
		err := cs.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, builder.String()))
		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		err := cs.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Недоступная команда"))
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
}
