package message_service

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/commands"
	context2 "TeleRequestHandler/controller/state_machine"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

type TelegramMessageService struct {
	commands map[string]commands.Command
	context  map[int64]*context2.UserContext
	bot      bot.Bot[any, string, int64]
}

func NewTelegramMessageService(commands map[string]commands.Command, context map[int64]*context2.UserContext, bot bot.Bot[any, string, int64]) *TelegramMessageService {
	return &TelegramMessageService{commands, context, bot}
}

func (ms *TelegramMessageService) ProcessMessages(ctx context.Context, updChan tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case upd := <-updChan:
			if upd.Message == nil {
				continue
			}
			usrCtx, ok := ms.context[upd.Message.Chat.ID]
			if !ok {
				usrCtx = context2.GetDefaultContext(upd.Message.Chat.ID)
				ms.context[upd.Message.Chat.ID] = usrCtx

			}
			currCmd, ok := ms.commands[usrCtx.CommandName]
			if !ok && upd.Message.IsCommand() {
				currCmd, ok = ms.commands[upd.Message.Text]
			}
			if ok {
				currCmd.Execute(usrCtx, upd)
				continue
			} else {
				err := ms.bot.SendMessage(usrCtx.ChatId, "Некорректная команда, введите /help, чтобы узнать о доступных командах")
				if err != nil {
					logger.Error(err.Error())
				}
				continue
			}
		}
	}
}
