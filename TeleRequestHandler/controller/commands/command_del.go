package commands

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/state_machine"
	"TeleRequestHandler/custom_erros"
	"TeleRequestHandler/repo_service"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandDelHandler struct {
	bot              bot.Bot[any, tgbotapi.MessageConfig]
	repoService      repo_service.RepoRegisterService
	stateTransitions map[state_machine.StateName]*state_machine.State
}

func NewCommandDelHandler(bot bot.Bot[any, tgbotapi.MessageConfig], repoService repo_service.RepoRegisterService) *CommandDelHandler {
	return &CommandDelHandler{bot, repoService, state_machine.GetTransitions("del", bot)}
}

func (ca *CommandDelHandler) Execute(usrCtx *state_machine.UserContext, update tgbotapi.Update) {
	if usrCtx.CurrentState.Name != state_machine.NONE {
		err := usrCtx.CurrentState.Process(usrCtx, update)
		if err != nil {
			logger.Error(err.Error())
			var prErr custom_erros.ProcessError
			if errors.As(err, &prErr) {
				if prErr.Error() != "" {
					err = ca.bot.Send(tgbotapi.NewMessage(usrCtx.ChatId, prErr.Error()))
					if err != nil {
						logger.Error(err.Error())
					}
				}
				err = usrCtx.CurrentState.Start(usrCtx)
				if err != nil {
					logger.Error(err.Error())
				}
				return
			}
		}
	}
	usrCtx.CurrentState, _ = ca.stateTransitions[usrCtx.CurrentState.Name]
	if usrCtx.CurrentState.Name == state_machine.NONE {
		if usrCtx.Link != "" {
			err := ca.repoService.DeleteRepo(usrCtx.ChatId, usrCtx.Link)
			var servErr custom_erros.ServerError
			if errors.As(err, &servErr) && servErr.StatusCode == 404 {
				err = ca.bot.Send(tgbotapi.NewMessage(usrCtx.ChatId, "Вы не отслеживаете данный репозиторий"))
				if err != nil {
					logger.Error(err.Error())
				}
				usrCtx.CurrentState = state_machine.NewState(state_machine.NONE, state_machine.NewNoneState(ca.bot))
			} else {
				err = ca.bot.Send(tgbotapi.NewMessage(usrCtx.ChatId, "Некорректная ссылка"))
				if err != nil {
					logger.Error(err.Error())
				}
			}
		} else {
			err := ca.bot.Send(tgbotapi.NewMessage(usrCtx.ChatId, "Ошибка при вводе данных, попробуйте заново"))
			if err != nil {
				logger.Error(err.Error())
			}
		}
		usrCtx.CommandName = ""
	} else {
		err := usrCtx.CurrentState.Start(usrCtx)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
