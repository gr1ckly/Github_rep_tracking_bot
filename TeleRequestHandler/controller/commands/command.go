package commands

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/chat_service"
	"TeleRequestHandler/repo_service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

type Command struct {
	name        string
	description string
	CommandHandler
}

func NewCommand(name string, description string, handler CommandHandler) *Command {
	return &Command{name, description, handler}
}

func (c *Command) GetDescription() string {
	return c.name + ": " + c.description
}

func GetCommands(bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig], chatService chat_service.ChatRegisterService, repoService repo_service.RepoRegisterService) map[string]*Command {
	return map[string]*Command{
		"help":  NewCommand("help", "Справка о доступных командах", NewCommandHelpHandler(bot)),
		"start": NewCommand("start", "Начало работы с ботом", NewCommandStartHandler(bot, chatService)),
		"add":   NewCommand("add", "Начало отслеживания нового репозитория", NewCommandAddHandler(bot, repoService)),
		"del":   NewCommand("del", "Прекращение отслеживания репозитория", NewCommandDelHandler(bot, repoService)),
		"repos": NewCommand("repos", "Вывод отслеживаемых репозиториев", NewCommandReposHandler(bot, repoService)),
	}
}
