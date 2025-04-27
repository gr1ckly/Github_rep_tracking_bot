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
	commands     map[string]commands.Command
	contextMutex sync.RWMutex
	context      map[int64]context2.UserContext
	bot          bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
	wg           sync.WaitGroup
}

func NewTelegramMessageService(commands map[string]commands.Command, context map[int64]context2.UserContext, bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) *TelegramMessageService {
	return &TelegramMessageService{sync.RWMutex{}, commands, sync.RWMutex{}, context, bot, sync.WaitGroup{}}
}

func (ms *TelegramMessageService) ProcessMessages(ctx context.Context, updChan tgbotapi.UpdatesChannel) {
	defer ms.wg.Wait()
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		case upd := <-updChan:
			if upd.CallbackQuery != nil {
				go func() {
					ms.wg.Add(1)
					defer ms.wg.Done()
					ms.contextMutex.RLock()
					currUser, ok := ms.context[upd.CallbackQuery.From.ID]
					ms.contextMutex.RUnlock()
					if !ok {
						currUser = context2.GetDefaultContext(upd.CallbackQuery.From.ID)
						currUser.CurrentState = context2.NewState(context2.NONE, context2.NewNoneState(ms.bot))
						currUser, err = currUser.CurrentState.Start(currUser)
						if err != nil {
							logger.Error(err.Error())
						}
					}
					currUser, err = ms.processMessage(upd, currUser)
					if err != nil {
						logger.Error(err.Error())
					}
					ms.contextMutex.Lock()
					ms.context[upd.CallbackQuery.From.ID] = currUser
					ms.contextMutex.Unlock()
				}()
			} else if upd.Message != nil {
				if upd.Message.IsCommand() {
					go func() {
						ms.wg.Add(1)
						defer ms.wg.Done()
						ms.contextMutex.RLock()
						currUser, ok := ms.context[upd.Message.From.ID]
						ms.contextMutex.RUnlock()
						if !ok {
							currUser = context2.GetDefaultContext(upd.Message.From.ID)
							currUser.CurrentState = context2.NewState(context2.NONE, context2.NewNoneState(ms.bot))
							currUser, err = currUser.CurrentState.Start(currUser)
							if err != nil {
								logger.Error(err.Error())
							}
						}
						currUser, err = ms.processCommand(upd, currUser)
						if err != nil {
							logger.Error(err.Error())
						}
						ms.contextMutex.Lock()
						ms.context[upd.Message.From.ID] = currUser
						ms.contextMutex.Unlock()
					}()
				} else {
					go func() {
						ms.wg.Add(1)
						defer ms.wg.Done()
						ms.contextMutex.RLock()
						currUser, ok := ms.context[upd.Message.From.ID]
						ms.contextMutex.RUnlock()
						if !ok {
							currUser = context2.GetDefaultContext(upd.Message.From.ID)
							currUser.CurrentState = context2.NewState(context2.NONE, context2.NewNoneState(ms.bot))
							currUser, err = currUser.CurrentState.Start(currUser)
							if err != nil {
								logger.Error(err.Error())
							}
						}
						currUser, err := ms.processMessage(upd, currUser)
						if err != nil {
							logger.Error(err.Error())
						}
						ms.contextMutex.Lock()
						ms.context[upd.Message.From.ID] = currUser
						ms.contextMutex.Unlock()
					}()
				}
			}
		}
	}
}

func (ms *TelegramMessageService) processMessage(update tgbotapi.Update, usrCtx context2.UserContext) (context2.UserContext, error) {
	if usrCtx.CommandName != "" {
		ms.commandMutex.RLock()
		currCmd := ms.commands[usrCtx.CommandName]
		ms.commandMutex.RUnlock()
		return currCmd.Execute(usrCtx, update), nil
	} else {
		return usrCtx, ms.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Введите /help, чтобы увидеть список доступных команд"))
	}
}

func (ms *TelegramMessageService) processCommand(update tgbotapi.Update, usrCtx context2.UserContext) (context2.UserContext, error) {
	if usrCtx.CommandName != "" && update.Message.Command() != "cancel" {
		return usrCtx, ms.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Чтобы отменить команду введите /cancel"))
	} else {
		ms.commandMutex.RLock()
		currCmd, ok := ms.commands[update.Message.Command()]
		ms.commandMutex.RUnlock()
		if ok {
			usrCtx.CommandName = update.Message.Command()
			return currCmd.Execute(usrCtx, update), nil
		} else {
			return usrCtx, ms.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Некорректная команда, введите /help, чтобы просмотреть список доступных команд"))
		}
	}
}
