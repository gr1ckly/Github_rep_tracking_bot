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
	bot         bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
	chatService chat_service.ChatRegisterService
}

func NewCommandStartHandler(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig], chatService chat_service.ChatRegisterService) CommandStartHandler {
	return CommandStartHandler{bot, chatService}
}

func (cs CommandStartHandler) Execute(usrCtx state_machine.UserContext, upd tgbotapi.Update) state_machine.UserContext {
	err := cs.chatService.RegisterChat(converters.ConvertChat(usrCtx, upd))
	if err == nil {
		usrCtx.CurrentState = state_machine.NewState(state_machine.NONE, state_machine.NewNoneState(cs.bot))
		err = cs.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Вы успешно зарегистрировались, введите /help, чтобы получить информацию о доступных командах"))
		if err != nil {
			logger.Error(err.Error())
		}
		usrCtx, err = usrCtx.CurrentState.Start(usrCtx)
		if err != nil {
			logger.Error(err.Error())
		}
		return usrCtx
	}
	var servErr custom_erros.ServerError
	if errors.As(err, &servErr) && servErr.StatusCode == 409 {
		err = cs.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Вы уже зарегистрированы"))
		if err != nil {
			logger.Error(err.Error())
		}
		usrCtx.CurrentState = state_machine.NewState(state_machine.NONE, state_machine.NewNoneState(cs.bot))
		usrCtx, err = usrCtx.CurrentState.Start(usrCtx)
		if err != nil {
			logger.Error(err.Error())
		}
	} else {
		err = cs.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Ошибка при регистрации чата, повторите попытку позднее"))
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return usrCtx
}
