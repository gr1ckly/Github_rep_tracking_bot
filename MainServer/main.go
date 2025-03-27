package main

import (
	"Crypto_Bot/MainServer/github_sdk"
	"Crypto_Bot/MainServer/server"
	"Crypto_Bot/MainServer/server/validators"
	"Crypto_Bot/MainServer/storage/postgres"
	"context"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
	serverUrl := os.Getenv("SERVER_HOST")
	if serverUrl == "" {
		return
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return
	}
	apiUrl := os.Getenv("GITHUB_URL")
	if apiUrl == "" {
		return
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return
	}
	acceptFormat := os.Getenv("GITHUB_ACCEPT_FORMAT")
	if acceptFormat == "" {
		return
	}
	apiVersion := os.Getenv("GITHUB_API_VERSION")
	if apiVersion == "" {
		return
	}
	timeout, err := strconv.Atoi(strings.TrimSpace(os.Getenv("GITHUB_TIMEOUT")))
	if err != nil {
		return
	}
	ghService := github_sdk.NewHttpGithubService(apiUrl, token, acceptFormat, apiVersion, timeout)
	chatStore, err := postgres.NewPostgresChatStore(context.Background(), dbUrl)
	if err != nil {
		return
	}
	repoStore, err := postgres.NewPostgresRepoStore(context.Background(), dbUrl)
	if err != nil {
		return
	}
	chatRepoRecordStore, err := postgres.NewPostgresChatRepoRecordStore(context.Background(), dbUrl)
	if err != nil {
		return
	}
	storeManager := server.NewStoreManager(chatStore, repoStore, chatRepoRecordStore)
	validator, err := validators.NewUrlValidator("^https:\\/\\/github\\.com\\/[a-zA-Z0-9_-]+\\/[a-zA-Z0-9_-]+(\\.git)?$", ghService)
	if err != nil {
		return
	}
	server := server.BuildServer(serverUrl, validator, storeManager)
	defer server.Stop()
	server.Start()
}
