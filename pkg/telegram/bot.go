package telegram

import (
	"Todoist-bot/pkg/storage/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
	"log"
)

type Bot struct {
	bot       *tgbotapi.BotAPI
	client    *todoist.Client
	db        *postgres.Postgres
	serverURL string
}

func NewBot(bot *tgbotapi.BotAPI, client *todoist.Client, db *postgres.Postgres, serverURL string) *Bot {
	return &Bot{bot: bot,
		client:    client,
		db:        db,
		serverURL: serverURL,
	}
}

func (b *Bot) Start() error {
	log.Printf("Starting %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()

	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				err := b.handleCommand(update.Message)
				if err != nil {
					return
				}
				continue
			} else {
				b.handleMessage(update.Message)
			}
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
