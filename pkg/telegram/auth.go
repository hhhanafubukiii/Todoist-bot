package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
	"log"
	"net/http"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	httpClient := &http.Client{}

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
