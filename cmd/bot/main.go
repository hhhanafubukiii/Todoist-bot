package main

import (
	postgres2 "Todoist-bot/pkg/storage/postgres"
	"Todoist-bot/pkg/telegram"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	clientId := os.Getenv("client_id")
	clientSecret := os.Getenv("client_secret")
	token := os.Getenv("token")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	httpClient := &http.Client{}
	db := &postgres2.Postgres{}

	client := todoist.NewClient(httpClient, clientId, clientSecret)
	telegramBot := telegram.NewBot(bot, client, db)
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
