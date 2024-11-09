package telegram

import (
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
	clientId            = os.Getenv("client_id")
	clientSecret        = os.Getenv("client_secret")
	scope               = os.Getenv("scope")
	secretString        = os.Getenv("secretString")
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
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please choose a command from the command list.\n/help")

	switch message.Command() {
	case commandStart:
		return b.handleCommandStart(message)
	default:
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) handleCommandStart(message *tgbotapi.Message) error {

}
