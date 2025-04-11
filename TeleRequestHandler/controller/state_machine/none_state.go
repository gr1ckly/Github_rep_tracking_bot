package state_machine

import (
	"TeleRequestHandler/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NoneState struct {
	bot bot.Bot[any, string, int64]
}

func NewNoneState(bot bot.Bot[any, string, int64]) *NoneState {
	return &NoneState{bot}
}

func (us *NoneState) Start(chatId int64) error {
	return us.bot.SendMessage(chatId, "Введите /help для получения информации о доступных командах")
}

func (us *NoneState) Process(usrCtx *UserContext, update tgbotapi.Update) error {
	return nil
}
