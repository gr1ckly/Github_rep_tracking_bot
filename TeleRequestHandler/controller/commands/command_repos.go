package commands

import (
	"Common"
	"TeleRequestHandler/bot"
	"TeleRequestHandler/controller/converters"
	"TeleRequestHandler/controller/state_machine"
	"TeleRequestHandler/repo_service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandReposHandler struct {
	bot              bot.Bot[any, string, int64]
	repoService      repo_service.RepoRegisterService
	stateTransitions map[state_machine.StateName]*state_machine.State
}

func NewCommandReposHandler(bot bot.Bot[any, string, int64], repoService repo_service.RepoRegisterService) *CommandReposHandler {
	return &CommandReposHandler{bot, repoService, state_machine.GetTransitions("repos", bot)}
}

func (cr *CommandReposHandler) Execute(usrCtx *state_machine.UserContext, update tgbotapi.Update) {
	if usrCtx.CurrentState.Name != state_machine.NONE {
		err := usrCtx.CurrentState.Process(usrCtx, update)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	usrCtx.CurrentState, _ = cr.stateTransitions[usrCtx.CurrentState.Name]
	if usrCtx.CurrentState.Name == state_machine.NONE {
		if usrCtx.Link != "" && usrCtx.Tags != nil {
			var repos []Common.RepoDTO
			if len(usrCtx.Tags) == 0 {
				newRepos, err := cr.repoService.GetReposByChat(usrCtx.ChatId)
				if err != nil {
					logger.Error(err.Error())
					err = cr.bot.SendMessage(usrCtx.ChatId, "Возникла ошибка при получении репозиториев")
					if err != nil {
						logger.Error(err.Error())
					}
					usrCtx.CommandName = ""
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
			err := cr.bot.SendMessage(usrCtx.ChatId, converters.ConvertToMessage(repos))
			if err != nil {
				logger.Error(err.Error())
			}
		} else {
			err := cr.bot.SendMessage(usrCtx.ChatId, "Ошибка при вводе данных, попробуйте заново")
			if err != nil {
				logger.Error(err.Error())
			}
		}
		usrCtx.CommandName = ""
	} else {
		err := usrCtx.CurrentState.Start(usrCtx.ChatId)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
