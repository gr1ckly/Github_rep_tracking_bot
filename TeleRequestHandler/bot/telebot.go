package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TeleBot struct {
	*tgbotapi.BotAPI
}

func NewTeleBot(token string) (*TeleBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	cmdDesc := GetCommandsDescription()
	commands := make([]tgbotapi.BotCommand, len(cmdDesc))
	pointer := 0
	for key, _ := range cmdDesc {
		commands[pointer] = tgbotapi.BotCommand{Command: key, Description: cmdDesc[key]}
		pointer++
	}
	config := tgbotapi.NewSetMyCommands(commands...)
	_, err = bot.Request(config)
	if err != nil {
		return nil, err
	}
	return &TeleBot{bot}, nil
}

func (bot *TeleBot) GetUpdates(timeout int) tgbotapi.UpdatesChannel {
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = timeout
	updates := bot.GetUpdatesChan(upd)
	return updates
}

func (bot *TeleBot) SendMessage(chatId int, msg string) error {
	return bot.SendMessage(chatId, msg)
}
