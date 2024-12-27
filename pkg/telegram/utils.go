package telegram

import (
	"fmt"
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

var token = os.Getenv("token")

func DeleteMessage(chatID int64, messageID int) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/deleteMessage?chat_id=%d&message_id=%d", token, chatID, messageID)
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func EditMessageReplyMarkup(chatID int64, messageID int, newKeyboard string) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageReplyMarkup?chat_id=%d&message_id=%d&reply_markup=%s",
		token,
		chatID,
		messageID,
		newKeyboard,
	)
	print(requestURL)
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func EditMessageText(chatID int64, messageID int, newText string) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText?chat_id=%d&message_id=%d&text=%s",
		token,
		chatID,
		messageID,
		newText,
	)
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func SendMessage(chatID int64, text, replyMarkup string) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s&reply_markup=%s",
		token,
		chatID,
		text,
		replyMarkup,
	)
	fmt.Println(requestURL)
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
