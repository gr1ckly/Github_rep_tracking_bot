package commands

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/converters"
	"TeleRequestHandler/controller/state_machine"
	"TeleRequestHandler/custom_erros"
	"TeleRequestHandler/repo_service"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandAddHandler struct {
	bot              bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
	repoService      repo_service.RepoRegisterService
	stateTransitions map[state_machine.StateName]state_machine.State
}

func NewCommandAddHandler(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig], repoService repo_service.RepoRegisterService) CommandAddHandler {
	return CommandAddHandler{bot, repoService, state_machine.GetTransitions("add", bot)}
}

func (ca CommandAddHandler) Execute(usrCtx state_machine.UserContext, update tgbotapi.Update) state_machine.UserContext {
	if usrCtx.CurrentState.Name != state_machine.NONE {
		usrCtx, err := usrCtx.CurrentState.Process(usrCtx, update)
		if err != nil {
			logger.Error(err.Error())
			var prErr custom_erros.ProcessError
			if errors.As(err, &prErr) {
				if prErr.Error() != "" {
					err = ca.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, prErr.Error()))
					if err != nil {
						logger.Error(err.Error())
					}
				}
				usrCtx, err = usrCtx.CurrentState.Start(usrCtx)
				if err != nil {
					logger.Error(err.Error())
				}
				return usrCtx
			}
		}
	}
	usrCtx.CurrentState, _ = ca.stateTransitions[usrCtx.CurrentState.Name]
	if usrCtx.CurrentState.Name == state_machine.NONE {
		if usrCtx.Link != "" && usrCtx.Events != nil {
			repo := converters.ConvertRepo(usrCtx)
			err := ca.repoService.AddRepo(usrCtx.ChatId, repo)
			var servErr custom_erros.ServerError
			if errors.As(err, &servErr) && servErr.StatusCode == 409 {
				if servErr.StatusCode == 409 {
					err = ca.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Вы уже отслеживаете данный репозиторий"))
					if err != nil {
						logger.Error(err.Error())
					}
					usrCtx.CurrentState = state_machine.NewState(state_machine.NONE, state_machine.NewNoneState(ca.bot))
				} else if servErr.StatusCode == 400 {
					err = ca.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Некорректная ссылка"))
					if err != nil {
						logger.Error(err.Error())
					}
				} else {
					err = ca.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Ошибка сервера"))
					if err != nil {
						logger.Error(err.Error())
					}
				}
				return usrCtx
			}
			if err == nil {
				err = ca.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Репозиторий удален из отслеживания"))
				if err != nil {
					logger.Error(err.Error())
				}
			} else {
				err = ca.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Ошибка выполнения программы"))
				if err != nil {
					logger.Error(err.Error())
				}
			}
		} else {
			err := ca.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Ошибка при вводе данных, попробуйте заново"))
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
	usrCtx, err := usrCtx.CurrentState.Start(usrCtx)
	if err != nil {
		logger.Error(err.Error())
	}
	return usrCtx
}
