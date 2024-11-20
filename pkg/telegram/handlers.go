package telegram

import (
	"Todoist-bot/pkg/config"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var (
	commandStart string = "start"
	databaseURL         = os.Getenv("databaseURL")
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	_, err := b.bot.Send(msg)
	if err != nil {
		log.Fatal("cannot send message!")
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, configs.Lexicon["response"]["unknown_command"])

	switch message.Command() {
	case commandStart:
		return b.handleCommandStart(message)
	default:
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) handleCommandStart(message *tgbotapi.Message) error {
	_, err := b.db.GetAccessToken(message.Chat.ID, databaseURL)
	if err != nil {
		b.initAuthorizationProcess(message)
	} else {
		msgText := configs.Lexicon["response"]["already_authorized"]
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)

		_, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
