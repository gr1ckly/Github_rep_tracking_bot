package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type TeleBot struct {
	*tgbotapi.BotAPI
}

func NewTeleBot(token string) (*TeleBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TeleBot{bot}, nil
}
