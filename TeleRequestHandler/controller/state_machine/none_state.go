package state_machine

import (
	"TeleRequestHandler/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NoneState struct {
	bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
}

func NewNoneState(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) *NoneState {
	return &NoneState{bot}
}

func (us *NoneState) Start(usrCtx UserContext) (UserContext, error) {
	usrCtx.CommandName = ""
	usrCtx.Tags = nil
	usrCtx.Events = nil
	usrCtx.Link = ""
	reply := tgbotapi.NewMessage(usrCtx.ChatId, "")
	reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	return usrCtx, us.bot.SendMessage(reply)
}

func (us *NoneState) Process(usrCtx UserContext, update tgbotapi.Update) (UserContext, error) {
	return usrCtx, nil
}
