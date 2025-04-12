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
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Launch(ctx context.Context) error {
	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		return fmt.Errorf("Couldn't find SERVER_URL")
	}
	chatApiExt := os.Getenv("CHAT_API_EXT")
	if chatApiExt == "" {
		return fmt.Errorf("Couldn't find CHAT_API_EXT")
	}
	serverTimeout, err := strconv.Atoi(os.Getenv("ANSWER_SERVER_TIMEOUT"))
	if err != nil {
		return err
	}
	chatService := chat_service.NewHttpChatRegisterService(serverUrl, chatApiExt, serverTimeout)
	repoApiExt := os.Getenv("REPO_API_EXT")
	if repoApiExt == "" {
		return fmt.Errorf("Couldn't find REPO_API_EXT")
	}
	repoService := repo_service.NewHttpRepoRegisterService(serverUrl, repoApiExt, serverTimeout)
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("Couldn't find BOT_TOKEN")
	}
	bot, err := bot2.NewTeleBot(botToken)
	if err != nil {
		return nil
	}
	kafkaNetwork := os.Getenv("KAFKA_NETWORK")
	if kafkaNetwork == "" {
		return fmt.Errorf("KAFKA_NETWORK not found")
	}
	kafkaAddr := os.Getenv("KAFKA_ADDR")
	if kafkaAddr == "" {
		return fmt.Errorf("KAFKA_ADDR not found")
	}
	kafkaTopicName := os.Getenv("KAFKA_TOPIC_NAME")
	if kafkaTopicName == "" {
		return fmt.Errorf("KAFKA_TOPIC_NAME")
	}
	kafkaTopicPartition, err := strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_TOPIC_PARTITION")))
	if err != nil {
		return fmt.Errorf("KAFKA_TOPIC_PARTITION incorrect")
	}
	kafkaTopicReplicationFactor, err := strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_TOPIC_REPLICATION_FACTOR")))
	if err != nil {
		return fmt.Errorf("KAFKA_TOPIC_REPLICATION_FACTOR incorrect")
	}
	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupId == "" {
		return fmt.Errorf("Couldn't find KAFKA_GROUP_ID")
	}
	kafkaService, err := notification_service.NewKafkaNotificationWaiter(kafkaNetwork, kafkaAddr, kafkaTopicName, kafkaTopicPartition, kafkaTopicReplicationFactor, kafkaGroupId, bot)
	defer kafkaService.Close()
	if err != nil {
		return err
	}
	telegramService := message_service.NewTelegramMessageService(commands.GetCommands(bot, chatService, repoService), map[int64]*state_machine.UserContext{}, bot)
	botTimeout, err := strconv.Atoi(os.Getenv("BOT_TIMEOUT"))
	if err != nil {
		return err
	}
	go func() {
		kafkaService.WaitNotification(ctx)
	}()
	telegramService.ProcessMessages(ctx, bot.GetUpdatesChannel(botTimeout))
	return nil
}
