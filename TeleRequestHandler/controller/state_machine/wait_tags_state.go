package state_machine

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/custom_erros"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type WaitTagsState struct {
	bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
}

func NewWaitTagsState(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) *WaitTagsState {
	return &WaitTagsState{bot}
}

func (wt *WaitTagsState) Start(usrCtx *UserContext) error {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "skip"),
		),
	)
	msg := tgbotapi.NewMessage(usrCtx.ChatId, "Введите теги через пробел")
	msg.ReplyMarkup = buttons
	return wt.bot.SendMessage(msg)
}

func (wt *WaitTagsState) Process(usrCtx *UserContext, update tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "skip":
			usrCtx.Tags = []string{}
			return nil
		default:
			return custom_erros.ProcessError{"Ошибка при вводе тегов, попробуйте заново"}
		}
	}
	usrCtx.Tags = strings.Split(update.Message.Text, " ")
	return wt.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Теги для данного репозитория: "+strings.Join(usrCtx.Tags, " ")))
}
