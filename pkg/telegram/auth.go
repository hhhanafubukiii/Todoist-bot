package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	httpClient := &http.Client{}

	var (
		clientId     = os.Getenv("client_id")
		clientSecret = os.Getenv("client_secret")
		secretString = os.Getenv("secretString")
		scope        = os.Getenv("scope")
	)

	client := todoist.NewClient(httpClient, clientId, clientSecret)
	authLink, err := client.GetAuthenticationURL(clientId, scope, secretString)
	if err != nil {
		log.Fatal("cannot get authorization URL!", err)
	}

	msgText := fmt.Sprintf("Привет! Чтобы я мог взаимодействовать с твоим аккаунтом Todoist, тебе необходимо дать мне на это доступ. Для этого переходи по ссылке:\n%s", authLink)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err = b.bot.Send(msg)
	if err != nil {
		log.Fatal("cannot send message!", err)
	}
	return nil
}

func GetAccessToken(chatId int64) (string, error) {
	return "", nil
}
