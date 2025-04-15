package message_service

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/commands"
	context2 "TeleRequestHandler/controller/state_machine"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
	"sync"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

type TelegramMessageService struct {
	commandMutex sync.RWMutex
	commands     map[string]*commands.Command
	contextMutex sync.RWMutex
	context      map[int64]*context2.UserContext
	bot          bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
}

func NewTelegramMessageService(commands map[string]*commands.Command, context map[int64]*context2.UserContext, bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) *TelegramMessageService {
	return &TelegramMessageService{sync.RWMutex{}, commands, sync.RWMutex{}, context, bot}
}

func (ms *TelegramMessageService) ProcessMessages(ctx context.Context, updChan tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case upd := <-updChan:
			var currUser int64
			if upd.Message != nil {
				currUser = upd.Message.From.ID
			} else if upd.CallbackQuery != nil {
				currUser = upd.CallbackQuery.From.ID
			} else {
				continue
			}
			usrCtx, ok := ms.context[currUser]
			if upd.Message != nil {
				if ok && upd.Message.Command() == "cancel" {
					usrCtx.CurrentState = context2.NewState(context2.NONE, context2.NewNoneState(ms.bot))
					err := usrCtx.CurrentState.Start(usrCtx)
					if err != nil {
						logger.Error(err.Error())
						continue
					}
				}
			}
			if !ok {
				usrCtx = context2.GetDefaultContext(currUser)
				usrCtx.CurrentState = context2.NewState(context2.NONE, context2.NewNoneState(ms.bot))
				ms.context[currUser] = usrCtx
			}
			currCmd, ok := ms.commands[usrCtx.CommandName]
			if !ok && upd.Message != nil && upd.Message.IsCommand() {
				currCmd, ok = ms.commands[upd.Message.Command()]
			}
			if ok {
				currCmd.Execute(usrCtx, upd)
				continue
			} else {
				err := ms.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Некорректная команда, введите /help, чтобы узнать о доступных командах"))
				if err != nil {
					logger.Error(err.Error())
				}
				continue
			}
		}
	}
}
