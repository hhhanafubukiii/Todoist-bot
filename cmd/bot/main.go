package main

import (
	"Todoist-bot/pkg/server"
	postgres "Todoist-bot/pkg/storage/postgres"
	"Todoist-bot/pkg/telegram"
	"github.com/looplab/fsm"
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

var (
	clientId     = os.Getenv("client_id")
	clientSecret = os.Getenv("client_secret")
	token        = os.Getenv("token")
	botURL       = os.Getenv("botURL")
	serverURL    = os.Getenv("serverURL")
)

func main() {
	httpClient := &http.Client{}
	db := &postgres.Postgres{}
	client := todoist.NewClient(httpClient, clientId, clientSecret)
	newFSM := fsm.NewFSM(
		"default",
		fsm.Events{
			{Name: "addTaskName", Src: []string{"default"}, Dst: "addTaskPriority"},
			{Name: "addTaskPriority", Src: []string{"addTaskName"}, Dst: "addTaskDeadline"},
			{Name: "addTaskDeadline", Src: []string{"addTaskPriority"}, Dst: "addTaskDescription"},
			{Name: "addTaskDescription", Src: []string{"addTaskDeadline"}, Dst: "default"},
		},
		fsm.Callbacks{},
	)

	go startBot(client, db, token, serverURL, newFSM)
	startServer(client, db, botURL)
}

func startBot(client *todoist.Client, db *postgres.Postgres, token, serverURL string, fsm *fsm.FSM) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, client, db, serverURL, fsm)

	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func startServer(client *todoist.Client, db *postgres.Postgres, botURL string) {
	authServer := server.NewAuthorizationServer(db, client, botURL)

	log.Println("Starting server")
	err := authServer.Start()
	if err != nil {
		log.Fatal("cannot start auth server", err)
	}
}
