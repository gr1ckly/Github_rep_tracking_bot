package commands

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/state_machine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandCancelHandler struct {
	bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
}

func NewCommandCancelHandler(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) CommandCancelHandler {
	return CommandCancelHandler{bot}
}

func (cc CommandCancelHandler) Execute(usrCtx state_machine.UserContext, upd tgbotapi.Update) state_machine.UserContext {
	if usrCtx.CommandName != "" {
		err := cc.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Отмена /"+usrCtx.CommandName))
		if err != nil {
			logger.Error(err.Error())
		}
	} else {
		err := cc.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Никакая команда не выполняется сейчас"))
		if err != nil {
			logger.Error(err.Error())
		}
	}
	usrCtx.CurrentState = state_machine.NewState(state_machine.NONE, state_machine.NewNoneState(cc.bot))
	usrCtx, err := usrCtx.CurrentState.Start(usrCtx)
	if err != nil {
		logger.Error(err.Error())
	}
	return usrCtx
}
