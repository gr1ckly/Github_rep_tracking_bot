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
	bot      bot.Bot[any, tgbotapi.MessageConfig]
}

func NewTelegramMessageService(commands map[string]commands.Command, context map[int64]*context2.UserContext, bot bot.Bot[any, tgbotapi.MessageConfig]) *TelegramMessageService {
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
			if ok && upd.Message.Command() == "cancel" {
				usrCtx.CurrentState = context2.NewState(context2.NONE, context2.NewNoneState(ms.bot))
				err := usrCtx.CurrentState.Start(usrCtx)
				if err != nil {
					logger.Error(err.Error())
					continue
				}
			}
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
				err := ms.bot.Send(tgbotapi.NewMessage(usrCtx.ChatId, "Некорректная команда, введите /help, чтобы узнать о доступных командах"))
				if err != nil {
					logger.Error(err.Error())
				}
				continue
			}
		}
	}
}
