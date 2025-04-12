package state_machine

import (
	"Common"
	"TeleRequestHandler/bot"
	"TeleRequestHandler/custom_erros"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"slices"
	"strings"
)

type WaitEventsState struct {
	bot            bot.Bot[any, tgbotapi.MessageConfig]
	checkboxStates map[int64]map[string]bool
}

func NewWaitEventsState(bot bot.Bot[any, tgbotapi.MessageConfig]) *WaitEventsState {
	return &WaitEventsState{bot, map[int64]map[string]bool{}}
}

func (we *WaitEventsState) Start(usrCtx *UserContext) error {
	reply := tgbotapi.NewMessage(usrCtx.ChatId, "")
	reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	state := we.checkboxStates[usrCtx.ChatId]
	if state == nil {
		state = make(map[string]bool)
		we.checkboxStates[usrCtx.ChatId] = state
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, opt := range Common.Events {
		checked := state[opt]
		icon := "⬜"
		if checked {
			icon = "☑️"
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %s", icon, opt), "event_"+opt)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	confirm := tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "confirm")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(confirm))
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return we.bot.Send(reply)
}

func (we *WaitEventsState) Process(usrCtx *UserContext, update tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, "event_") {
			event := strings.TrimPrefix(update.CallbackQuery.Data, "event_")
			if !slices.Contains(usrCtx.Events, event) {
				usrCtx.Events = append(usrCtx.Events, event)
			}
			return custom_erros.ProcessError{""}
		} else if update.CallbackQuery.Data == "confirm" {
			if usrCtx.Events == nil || len(usrCtx.Events) == 0 {
				return custom_erros.ProcessError{"Вы не выбрали ни одного события для отслеживания"}
			}
		}
	}
	return nil
}
