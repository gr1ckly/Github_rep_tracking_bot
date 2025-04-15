package state_machine

import (
	"Common"
	"TeleRequestHandler/bot"
	"TeleRequestHandler/custom_erros"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"slices"
	"strings"
)

type WaitEventsState struct {
	bot            bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
	checkboxStates map[int64]map[string]bool
}

func NewWaitEventsState(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) *WaitEventsState {
	return &WaitEventsState{bot, map[int64]map[string]bool{}}
}

func (we *WaitEventsState) Start(usrCtx *UserContext) error {
	builder := strings.Builder{}
	builder.WriteString("Выберите тип отслеживаемого события через пробел ( ")
	for _, event := range Common.Events {
		builder.WriteString(event)
		builder.WriteString(" ")
	}
	builder.WriteRune(')')
	return we.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, builder.String()))
}

func (we *WaitEventsState) Process(usrCtx *UserContext, update tgbotapi.Update) error {
	events := strings.Split(update.Message.Text, " ")
	for _, evnt := range events {
		if slices.Contains(Common.Events, evnt) && !slices.Contains(usrCtx.Events, evnt) {
			usrCtx.Events = append(usrCtx.Events, evnt)
		}
	}
	if usrCtx.Events == nil || len(usrCtx.Events) == 0 {
		return custom_erros.ProcessError{"Вы не выбрали ни одного события, введите данные еще раз"}
	}
	return nil
}
