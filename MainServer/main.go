package main

import (
	"Crypto_Bot/MainServer/LinkTracker"
	"Crypto_Bot/MainServer/github_sdk"
	"Crypto_Bot/MainServer/server"
	"Crypto_Bot/MainServer/server/validators"
	"Crypto_Bot/MainServer/storage/postgres"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	serverUrl := os.Getenv("MAIN_SERVER_HOST")
	if serverUrl == "" {
		logger.Error("MAIN_SERVER_HOST not found")
		return
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		logger.Error("DB_URL not found")
		return
	}
	apiUrl := os.Getenv("GITHUB_URL")
	if apiUrl == "" {
		logger.Error("GITHUB_URL not found")
		return
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		logger.Error("GITHUB_TOKEN not found")
		return
	}
	acceptFormat := os.Getenv("GITHUB_ACCEPT_FORMAT")
	if acceptFormat == "" {
		logger.Error("GITHUB_ACCEPT_FORMAT not found")
		return
	}
	apiVersion := os.Getenv("GITHUB_API_VERSION")
	if apiVersion == "" {
		logger.Error("GITHUB_API_VERSION not found")
		return
	}
	ghTimeout, err := strconv.Atoi(strings.TrimSpace(os.Getenv("GITHUB_TIMEOUT")))
	if err != nil {
		logger.Error("GITHUB_TIMEOUT incorrect")
		return
	}
	dbTimeout, err := strconv.Atoi(strings.TrimSpace(os.Getenv("DB_TIMEOUT")))
	if err != nil {
		logger.Error("DB_TIMEOUT incorrect")
		return
	}
	kafkaTimeout, err := strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_TIMEOUT")))
	if err != nil {
		logger.Error("KAFKA_TIMEOUT incorrect")
		return
	}
	batchSize, err := strconv.Atoi(strings.TrimSpace(os.Getenv("BATCH_SIZE")))
	if err != nil {
		logger.Error("BATCH_SIZE incorrect")
		return
	}
	ghService := github_sdk.NewHttpGithubService(apiUrl, token, acceptFormat, apiVersion, ghTimeout)
	chatStore, err := postgres.NewPostgresChatStore(dbTimeout, dbUrl)
	defer chatStore.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	repoStore, err := postgres.NewPostgresRepoStore(dbTimeout, dbUrl)
	defer repoStore.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	chatRepoRecordStore, err := postgres.NewPostgresChatRepoRecordStore(dbTimeout, dbUrl)
	defer chatRepoRecordStore.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	storeManager := server.NewStoreManager(chatStore, repoStore, chatRepoRecordStore)
	validator, err := validators.NewUrlValidator("^https:\\/\\/github\\.com\\/[a-zA-Z0-9_-]+\\/[a-zA-Z0-9_-]+(\\.git)?$", ghService)
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
		logger.Error("KAFKA_TOPIC_NAME not found")
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
	kafkaNotificationManager, err := LinkTracker.NewNotificationService(kafkaTimeout, kafkaNetwork, kafkaAddr, kafkaTopicName, kafkaTopicPartition, kafkaTopicReplicationFactor)
	linkTracker, err := LinkTracker.NewLinkTracker(ghService, storeManager, chatRepoRecordStore, batchSize)
	if err != nil {
		logger.Error(err.Error())
	}
	linkTracker.AddObserver(kafkaNotificationManager)
	linkTracker.StartTracking()
	defer linkTracker.Stop()
	server := server.BuildServer(serverUrl, validator, storeManager)
	defer server.Stop()
	err = server.Start()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
