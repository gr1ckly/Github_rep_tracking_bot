package commands

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/chat_service"
	"TeleRequestHandler/controller/converters"
	"TeleRequestHandler/controller/state_machine"
	"TeleRequestHandler/custom_erros"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandStartHandler struct {
	bot         bot.Bot[any, string, int64]
	chatService chat_service.ChatRegisterService
}

func NewCommandStartHandler(bot bot.Bot[any, string, int64], chatService chat_service.ChatRegisterService) *CommandStartHandler {
	return &CommandStartHandler{bot, chatService}
}

func (cs *CommandStartHandler) Execute(usrCtx *state_machine.UserContext, upd tgbotapi.Update) {
	err := cs.chatService.RegisterChat(converters.ConvertChat(usrCtx, upd))
	if err == nil {
		usrCtx.CurrentState = state_machine.NewState(state_machine.NONE, state_machine.NewNoneState(cs.bot))
		err = cs.bot.SendMessage(usrCtx.ChatId, "Вы успешно зарегистрировались, введите /help, чтобы получить информацию о доступных командах")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}
	var servErr custom_erros.ServerError
	if errors.As(err, &servErr) && servErr.StatusCode == 409 {
		err = cs.bot.SendMessage(usrCtx.ChatId, "Вы уже зарегистрированы")
		if err != nil {
			logger.Error(err.Error())
		}
		usrCtx.CurrentState = state_machine.NewState(state_machine.NONE, state_machine.NewNoneState(cs.bot))
	} else {
		err = cs.bot.SendMessage(usrCtx.ChatId, "Ошибка при регистрации чата, повторите попытку позднее")
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
