package main

import (
	"TeleRequestHandler/telebot_app"
	"context"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = telebot_app.Launch(ctx)
	if err != nil {
		logger.Error(err.Error())
	}
}
