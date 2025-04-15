package commands

import (
	"Common"
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/converters"
	"TeleRequestHandler/controller/state_machine"
	"TeleRequestHandler/custom_erros"
	"TeleRequestHandler/repo_service"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandReposHandler struct {
	bot              bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]
	repoService      repo_service.RepoRegisterService
	stateTransitions map[state_machine.StateName]*state_machine.State
}

func NewCommandReposHandler(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig], repoService repo_service.RepoRegisterService) *CommandReposHandler {
	return &CommandReposHandler{bot, repoService, state_machine.GetTransitions("repos", bot)}
}

func (cr *CommandReposHandler) Execute(usrCtx *state_machine.UserContext, update tgbotapi.Update) {
	if usrCtx.CurrentState.Name != state_machine.NONE {
		err := usrCtx.CurrentState.Process(usrCtx, update)
		if err != nil {
			logger.Error(err.Error())
			var prErr custom_erros.ProcessError
			if errors.As(err, &prErr) {
				if prErr.Error() != "" {
					err = cr.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, prErr.Error()))
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
	} else {
		usrCtx.CommandName = "repos"
	}
	usrCtx.CurrentState, _ = cr.stateTransitions[usrCtx.CurrentState.Name]
	if usrCtx.CurrentState.Name == state_machine.NONE {
		var repos []Common.RepoDTO
		if len(usrCtx.Tags) == 0 {
			newRepos, err := cr.repoService.GetReposByChat(usrCtx.ChatId)
			if err != nil {
				logger.Error(err.Error())
				err = cr.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, "Возникла ошибка при получении репозиториев"))
				if err != nil {
					logger.Error(err.Error())
				}
				return
			}
			repos = append(repos, newRepos...)
		} else {
			for _, tag := range usrCtx.Tags {
				newRepos, err := cr.repoService.GetReposByTag(usrCtx.ChatId, tag)
				if err != nil {
					logger.Error(err.Error())
					continue
				}
				repos = append(repos, newRepos...)
			}
		}
		err := cr.bot.SendMessage(tgbotapi.NewMessage(usrCtx.ChatId, converters.ConvertToMessage(repos)))
		if err != nil {
			logger.Error(err.Error())
		}
	}
	err := usrCtx.CurrentState.Start(usrCtx)
	if err != nil {
		logger.Error(err.Error())
	}
}
