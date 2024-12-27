package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

// aaa
import (
	"Todoist-bot/pkg/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var (
	clientId     = os.Getenv("client_id")
	clientSecret = os.Getenv("client_secret")
	scope        = os.Getenv("scope")
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) {
	httpClient := &http.Client{}
	client := todoist.NewClient(httpClient, clientId, clientSecret)

	authLink, err := client.GetAuthenticationURL(clientId, scope, strconv.FormatInt(message.Chat.ID, 10))
	if err != nil {
		log.Fatal("cannot get authorization URL!", err)
	}

	msgText := fmt.Sprintf(configs.Lexicon["response"]["start"], authLink)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)

	_, err = b.bot.Send(msg)
	if err != nil {
		log.Fatal("cannot send message!", err)
	}
}
