package main

import (
	"TeleRequestHandler/telebot_app"
	"context"
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	telebot_app.Launch(ctx)
}
