package telebot_app

import (
	bot2 "TeleRequestHandler/bot"
	"TeleRequestHandler/chat_service"
	"TeleRequestHandler/controller/commands"
	"TeleRequestHandler/controller/message_service"
	"TeleRequestHandler/controller/state_machine"
	"TeleRequestHandler/notification_service"
	"TeleRequestHandler/repo_service"
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

func Launch(ctx context.Context) {
	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		logger.Error("Couldn't find SERVER_URL")
		return
	}
	chatApiExt := os.Getenv("CHAT_API_EXT")
	if chatApiExt == "" {
		logger.Error("Couldn't find CHAT_API_EXT")
		return
	}
	serverTimeout, err := strconv.Atoi(os.Getenv("ANSWER_SERVER_TIMEOUT"))
	if err != nil {
		logger.Error(err.Error())
		return
	}
	chatService := chat_service.NewHttpChatRegisterService(serverUrl, chatApiExt, serverTimeout)
	repoApiExt := os.Getenv("REPO_API_EXT")
	if repoApiExt == "" {
		logger.Error("Couldn't find REPO_API_EXT")
		return
	}
	repoService := repo_service.NewHttpRepoRegisterService(serverUrl, repoApiExt, serverTimeout)
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		logger.Error("Couldn't find BOT_TOKEN")
		return
	}
	bot, err := bot2.NewTeleBot(botToken)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	kafkaNetwork := os.Getenv("KAFKA_NETWORK")
	if kafkaNetwork == "" {
		logger.Error("KAFKA_NETWORK not found")
		return
	}
	kafkaAddr := os.Getenv("KAFKA_ADDR")
	if kafkaAddr == "" {
		logger.Error("KAFKA_ADDR not found")
		return
	}
	kafkaTopicName := os.Getenv("KAFKA_TOPIC_NAME")
	if kafkaTopicName == "" {
		logger.Error("KAFKA_TOPIC_NAME")
		return
	}
	kafkaTopicPartition, err := strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_TOPIC_PARTITION")))
	if err != nil {
		logger.Error("KAFKA_TOPIC_PARTITION incorrect")
		return
	}
	kafkaTopicReplicationFactor, err := strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_TOPIC_REPLICATION_FACTOR")))
	if err != nil {
		logger.Error("KAFKA_TOPIC_REPLICATION_FACTOR incorrect")
		return
	}
	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupId == "" {
		logger.Error("Couldn't find KAFKA_GROUP_ID")
		return
	}
	kafkaService, err := notification_service.NewKafkaNotificationWaiter(kafkaNetwork, kafkaAddr, kafkaTopicName, kafkaTopicPartition, kafkaTopicReplicationFactor, kafkaGroupId, bot)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer kafkaService.Close()
	telegramService := message_service.NewTelegramMessageService(commands.GetCommands(bot, chatService, repoService), map[int64]state_machine.UserContext{}, bot)
	botTimeout, err := strconv.Atoi(os.Getenv("BOT_TIMEOUT"))
	if err != nil {
		logger.Error(err.Error())
		return
	}
	go func() {
		kafkaService.WaitNotification(ctx)
	}()
	telegramService.ProcessMessages(ctx, bot.GetUpdatesChannel(botTimeout))
}
