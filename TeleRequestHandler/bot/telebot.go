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
	return &TeleBot{bot}, nil
}

func (bot *TeleBot) GetUpdates(timeout int) tgbotapi.UpdatesChannel {
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60
	updates := bot.GetUpdatesChan(upd)
	return updates
}

func (bot *TeleBot) SendMessage(chatId int, msg string) error {
	return bot.SendMessage(chatId, msg)
}
