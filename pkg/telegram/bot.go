package telegram

import (
	configs "Todoist-bot/pkg/config"
	"Todoist-bot/pkg/storage/postgres"
	"Todoist-bot/pkg/storage/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
	"github.com/looplab/fsm"
	"log"
	"strconv"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	client   *todoist.Client
	postgres *postgres.Postgres
	redis    *redis.Redis
	fsm      *fsm.FSM
}

func NewBot(bot *tgbotapi.BotAPI, client *todoist.Client, db *postgres.Postgres, redis *redis.Redis, fsm *fsm.FSM) *Bot {
	return &Bot{bot: bot,
		client:   client,
		postgres: db,
		redis:    redis,
		fsm:      fsm,
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
	task := todoist.AddTask{}
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				if err := b.handleCommand(update.Message); err != nil {
					return
				}
				continue
			} else {
				switch b.fsm.Current() {
				case "addTaskName":
					if update.Message.Text == configs.Lexicon["buttons"]["cancel"] {
						if err := b.handleCancel(update.Message); err != nil {
							log.Fatal(err)
						}
					} else {
						task.Name = update.Message.Text
						if err := b.handleTaskName(update.Message); err != nil {
							log.Fatal(err)
						}
					}
				case "addTaskPriority":
					if update.Message.Text == configs.Lexicon["buttons"]["skip"] {
						task.Priority = "1"
						if err := b.handleTaskPriority(update.Message); err != nil {
							log.Fatal(err)
						}
					} else if update.Message.Text == configs.Lexicon["buttons"]["cancel"] {
						if err := b.handleCancel(update.Message); err != nil {
							log.Fatal(err)
						}
					} else {
						task.Priority = strconv.Itoa(int(update.Message.Text[0]))
						if err := b.handleTaskPriority(update.Message); err != nil {
							log.Fatal(err)
						}
					}
				case "addTaskDeadline":
					if update.Message.Text == configs.Lexicon["buttons"]["skip"] {
						task.DueDate = "today"
						if err := b.handleTaskDeadline(update.Message); err != nil {
							log.Fatal(err)
						}
					} else if update.Message.Text == configs.Lexicon["buttons"]["cancel"] {
						if err := b.handleCancel(update.Message); err != nil {
							log.Fatal(err)
						}
					} else {
						task.DueDate = update.Message.Text
						if err := b.handleTaskDeadline(update.Message); err != nil {
							log.Fatal(err)
						}
					}
				case "addTaskDescription":
					if update.Message.Text == configs.Lexicon["buttons"]["skip"] {
						task.Description = ""
						if err := b.handleTaskDescription(update.Message.Chat.ID, task.Name, task.Priority, task.DueDate, task.Description); err != nil {
							log.Fatal(err)
						}
					} else if update.Message.Text == configs.Lexicon["buttons"]["cancel"] {
						if err := b.handleCancel(update.Message); err != nil {
							log.Fatal(err)
						}
					} else {
						task.Description = update.Message.Text
						if err := b.handleTaskDescription(update.Message.Chat.ID, task.Name, task.Priority, task.DueDate, task.Description); err != nil {
							log.Fatal(err)
						}
					}
				case "updateTaskName":
					err := b.updateTaskName(update.Message)
					if err != nil {
						log.Fatal(err)
					}
				case "updateTodayTaskName":
					err := b.updateTaskName(update.Message)
					if err != nil {
						log.Fatal(err)
					}
				case "updateTaskPriority":
					err := b.updateTaskPriority(update.Message)
					if err != nil {
						log.Fatal(err)
					}
				case "updateTodayTaskPriority":
					err := b.updateTaskPriority(update.Message)
					if err != nil {
						log.Fatal(err)
					}
				case "updateTaskDueDate":
					err := b.updateTaskDueDate(update.Message)
					if err != nil {
						log.Fatal(err)
					}
				case "updateTodayTaskDueDate":
					err := b.updateTaskDueDate(update.Message)
					if err != nil {
						log.Fatal(err)
					}
				case "updateTaskDescription":
					err := b.updateTaskDescription(update.Message)
					if err != nil {
						log.Fatal(err)
					}
				case "updateTodayTaskDescription":
					err := b.updateTaskDescription(update.Message)
					if err != nil {
						log.Fatal(err)
					}
				default:
					switch update.Message.Text {
					case configs.Lexicon["buttons"]["new_task"]:
						err := b.handleAddTask(update.Message)
						if err != nil {
							log.Fatal(err)
						}
					case configs.Lexicon["buttons"]["today_tasks"]:
						if err := b.handleTodayTasksButton(update.Message); err != nil {
							log.Fatal(err)
						}
					case configs.Lexicon["buttons"]["all_tasks"]:
						if err := b.handleAllTasksButton(update.Message); err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		} else if update.CallbackQuery != nil {
			switch b.fsm.Current() {
			case "selectTask":
				if update.CallbackQuery.Data == configs.Lexicon["buttons"]["cancel"] {
					if err := b.handleInlineCancel(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else {
					if err := b.handleSelectTask(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				}
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := b.bot.AnswerCallbackQuery(callback); err != nil {
					log.Fatal(err)
				}
			case "selectTodayTask":
				if update.CallbackQuery.Data == configs.Lexicon["buttons"]["cancel"] {
					if err := b.handleInlineCancel(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else {
					if err := b.handleSelectTodayTask(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				}
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := b.bot.AnswerCallbackQuery(callback); err != nil {
					log.Fatal(err)
				}
			case "selectedTask":
				if update.CallbackQuery.Data == configs.Lexicon["buttons"]["back"] {
					if err := b.handleAllTasksCallBack(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["complete_task"] {
					if err := b.handleCompleteTask(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["update_task"] {
					if err := b.handleUpdateTask(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["delete_task"] {
					if err := b.handleDeleteTask(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				}
			case "selectedTodayTask":
				if update.CallbackQuery.Data == configs.Lexicon["buttons"]["back"] {
					if err := b.handleTodayTasksCallBack(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["complete_task"] {
					if err := b.handleCompleteTodayTask(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["update_task"] {
					if err := b.handleUpdateTodayTask(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["delete_task"] {
					if err := b.handleDeleteTask(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				}
			case "updateTask":
				if update.CallbackQuery.Data == configs.Lexicon["buttons"]["back"] {
					if err := b.handleBack(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["name"] {
					if err := b.handleUpdateTaskName(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["priority"] {
					if err := b.handleUpdateTaskPriority(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else if update.CallbackQuery.Data == configs.Lexicon["buttons"]["due_date"] {
					if err := b.handleUpdateTaskDueDate(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				} else {
					if err := b.handleUpdateTaskDescription(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				}
			case "updateTodayTask":
				if update.CallbackQuery.Data == configs.Lexicon["buttons"]["back"] {
					if err := b.handleBack(update.CallbackQuery); err != nil {
						log.Fatal(err)
					}
				}
			}
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.bot.AnswerCallbackQuery(callback); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
