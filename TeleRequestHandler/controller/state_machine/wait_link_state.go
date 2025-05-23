package state_machine

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/custom_erros"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WaitLinkState struct {
	bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
}

func NewWaitLinkState(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) *WaitLinkState {
	return &WaitLinkState{bot}
}

func (wl *WaitLinkState) Start(usrCtx UserContext) (UserContext, error) {
	reply := tgbotapi.NewMessage(usrCtx.ChatId, "Введите ссылку на репозиторий")
	reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return usrCtx, wl.bot.SendMessage(reply)
}

func (wl *WaitLinkState) Process(usrCtx UserContext, update tgbotapi.Update) (UserContext, error) {
	usrCtx.Link = update.Message.Text
	if update.Message.Text == "" {
		return usrCtx, custom_erros.ProcessError{"Ошибка при вводе ссылки, попробуйте заново"}
	}
	return usrCtx, wl.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Текущая ссылка: "+usrCtx.Link))
}
