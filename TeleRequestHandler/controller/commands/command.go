package commands

import (
	"TeleRequestHandler/bot"
	"TeleRequestHandler/chat_service"
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

func GetCommands(bot bot.Bot[any, string, int64], chatService chat_service.ChatRegisterService) map[string]*Command {
	return map[string]*Command{
		"help":  NewCommand("help", "Справка о доступных командах", NewCommandHelpHandler(bot)),
		"start": NewCommand("start", "Начало работы с ботом", NewCommandStartHandler(bot, chatService)),
		"add":   NewCommand("add", "Начало отслеживания нового репозитория", nil),
		"del":   NewCommand("del", "Прекращение отслеживания репозитория", nil),
		"repos": NewCommand("repos", "Вывод отслеживаемых репозиториев", nil),
	}
}
