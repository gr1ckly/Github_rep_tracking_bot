package state_machine

import (
	"TeleRequestHandler/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NoneState struct {
	bot bot.Bot[any, tgbotapi.MessageConfig]
}

func NewNoneState(bot bot.Bot[any, tgbotapi.MessageConfig]) *NoneState {
	return &NoneState{bot}
}

func (us *NoneState) Start(usrCtx *UserContext) error {
	usrCtx.CommandName = ""
	usrCtx.Tags = nil
	usrCtx.Events = nil
	usrCtx.Link = ""
	reply := tgbotapi.NewMessage(usrCtx.ChatId, "")
	reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	reply.DisableNotification = true
	return us.bot.Send(reply)
}

func (us *NoneState) Process(usrCtx *UserContext, update tgbotapi.Update) error {
	return nil
}
