package telegram

import (
	configs "Todoist-bot/pkg/config"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var (
	commandStart string = "start"
	commandHelp  string = "help"
	databaseURL  string = os.Getenv("databaseURL")
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleCommandStart(message)
	case commandHelp:
		return b.handleCommandHelp(message)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, configs.Lexicon["response"]["unknown_command"])
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) handleCommandStart(message *tgbotapi.Message) error {
	_, err := b.postgres.Get(message.Chat.ID, databaseURL)
	if err != nil {
		b.initAuthorizationProcess(message)
	} else {
		msgText := configs.Lexicon["response"]["already_authorized"]
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		msg.ReplyMarkup = MainKeyboard
		_, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleCommandHelp(message *tgbotapi.Message) error {
	msgText := configs.Lexicon["response"]["help"]
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = MainKeyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleAddTask(message *tgbotapi.Message) error {
	msgText := configs.Lexicon["response"]["new_task"]
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = CancelKeyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	if err = b.fsm.Event(context.Background(), "addTaskName"); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleTaskName(message *tgbotapi.Message) error {
	msgText := configs.Lexicon["response"]["enter_task_priority"]
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = PriorityKeyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	if err = b.fsm.Event(context.Background(), "addTaskPriority"); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleTaskPriority(message *tgbotapi.Message) error {
	msgText := configs.Lexicon["response"]["enter_task_deadline"]
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = DueDateKeyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	if err = b.fsm.Event(context.Background(), "addTaskDeadline"); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleTaskDeadline(message *tgbotapi.Message) error {
	msgText := configs.Lexicon["response"]["enter_task_description"]
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = CancelSkipKeyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	if err = b.fsm.Event(context.Background(), "addTaskDescription"); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleTaskDescription(chatId int64, name, priority, dueDate, description string) error {
	accessToken, err := b.postgres.Get(chatId, databaseURL)
	if err != nil {
		return err
	}

	if err = b.client.AddTask(name, priority, dueDate, description, accessToken); err != nil {
		return err
	}

	msgText := configs.Lexicon["response"]["add_task"]
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ReplyMarkup = MainKeyboard
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	if err = b.fsm.Event(context.Background(), "root"); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCancel(message *tgbotapi.Message) error {
	msgText := "OK"
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = MainKeyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "root")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleInlineCancel(callBack *tgbotapi.CallbackQuery) error {
	if err := DeleteMessage(callBack.Message.Chat.ID, callBack.Message.MessageID); err != nil {
		log.Fatal(err)
	}

	err := b.fsm.Event(context.Background(), "root")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleTodayTasksButton(message *tgbotapi.Message) error {
	accessToken, err := b.postgres.Get(message.Chat.ID, databaseURL)
	taskList, err := b.client.GetTodayTasks(accessToken)
	if err != nil {
		return err
	}

	tasksKeyboard := PrepareTaskListKeyboard(taskList)
	text := configs.Lexicon["response"]["task_list"]

	err = SendMessage(message.Chat.ID, text, tasksKeyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectTodayTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleTodayTasksCallBack(callBack *tgbotapi.CallbackQuery) error {
	accessToken, err := b.postgres.Get(callBack.Message.Chat.ID, databaseURL)
	taskList, err := b.client.GetTodayTasks(accessToken)
	if err != nil {
		return err
	}

	tasksKeyboard := PrepareTaskListKeyboard(taskList)
	text := configs.Lexicon["response"]["task_list"]

	err = EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, tasksKeyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectTodayTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleAllTasksButton(message *tgbotapi.Message) error {
	accessToken, err := b.postgres.Get(message.Chat.ID, databaseURL)
	taskList, err := b.client.GetAllTasks(accessToken)
	if err != nil {
		return err
	}

	tasksKeyboard := PrepareTaskListKeyboard(taskList)
	text := configs.Lexicon["response"]["task_list"]

	err = SendMessage(message.Chat.ID, text, tasksKeyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleAllTasksCallBack(callBack *tgbotapi.CallbackQuery) error {
	accessToken, err := b.postgres.Get(callBack.Message.Chat.ID, databaseURL)
	taskList, err := b.client.GetAllTasks(accessToken)
	if err != nil {
		return err
	}

	tasksKeyboard := PrepareTaskListKeyboard(taskList)
	text := configs.Lexicon["response"]["task_list"]

	err = EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, tasksKeyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleSelectTask(callBack *tgbotapi.CallbackQuery) error {
	taskID := callBack.Data
	chatID := callBack.Message.Chat.ID
	accessToken, err := b.postgres.Get(chatID, databaseURL)
	if err != nil {
		return err
	}

	err = b.redis.Save(strconv.FormatInt(chatID, 10), taskID)
	if err != nil {
		return err
	}

	task, err := b.client.GetTask(taskID, accessToken)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(configs.Lexicon["response"]["send_task_info"],
		task.Name,
		task.Priority,
		task.DueDate.DueString,
		task.Description,
	)
	keyboard := PrepareSelectedTaskKeyboard()

	err = EditMessageText(chatID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = EditMessageReplyMarkup(chatID, callBack.Message.MessageID, keyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectedTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleSelectTodayTask(callBack *tgbotapi.CallbackQuery) error {
	taskID := callBack.Data
	chatID := callBack.Message.Chat.ID
	accessToken, err := b.postgres.Get(chatID, databaseURL)
	if err != nil {
		return err
	}

	err = b.redis.Save(strconv.FormatInt(chatID, 10), taskID)
	if err != nil {
		return err
	}

	task, err := b.client.GetTask(taskID, accessToken)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(configs.Lexicon["response"]["send_task_info"],
		task.Name,
		task.Priority,
		task.DueDate.DueString,
		task.Description,
	)
	keyboard := PrepareSelectedTaskKeyboard()

	err = EditMessageText(chatID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = EditMessageReplyMarkup(chatID, callBack.Message.MessageID, keyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectedTodayTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCompleteTask(callBack *tgbotapi.CallbackQuery) error {
	accessToken, err := b.postgres.Get(callBack.Message.Chat.ID, databaseURL)
	if err != nil {
		return err
	}

	err, taskID := b.redis.Get(strconv.FormatInt(callBack.Message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.CloseTask(taskID, accessToken)
	if err != nil {
		return err
	}

	taskList, err := b.client.GetAllTasks(accessToken)
	if err != nil {
		return err
	}

	tasksKeyboard := PrepareTaskListKeyboard(taskList)
	text := configs.Lexicon["response"]["task_list"]

	err = EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, tasksKeyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCompleteTodayTask(callBack *tgbotapi.CallbackQuery) error {
	accessToken, err := b.postgres.Get(callBack.Message.Chat.ID, databaseURL)
	if err != nil {
		return err
	}

	err, taskID := b.redis.Get(strconv.FormatInt(callBack.Message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.CloseTask(taskID, accessToken)
	if err != nil {
		return err
	}

	taskList, err := b.client.GetTodayTasks(accessToken)
	if err != nil {
		return err
	}

	tasksKeyboard := PrepareTaskListKeyboard(taskList)
	text := configs.Lexicon["response"]["task_list"]

	err = EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, tasksKeyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectTodayTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTask(callBack *tgbotapi.CallbackQuery) error {
	keyboard := PrepareUpdateTaskKeyboard()

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, configs.Lexicon["response"]["update_task"])
	if err != nil {
		return err
	}
	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, keyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTodayTask(callBack *tgbotapi.CallbackQuery) error {
	keyboard := PrepareUpdateTaskKeyboard()

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, configs.Lexicon["response"]["update_task"])
	if err != nil {
		return err
	}
	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, keyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTodayTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleBack(callBack *tgbotapi.CallbackQuery) error {
	accessToken, _ := b.postgres.Get(callBack.Message.Chat.ID, databaseURL)
	err, taskID := b.redis.Get(strconv.FormatInt(callBack.Message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}
	task, err := b.client.GetTask(taskID, accessToken)
	if err != nil {
		return err
	}
	text := fmt.Sprintf(configs.Lexicon["response"]["send_task_info"],
		task.Name,
		task.Priority,
		task.DueDate.DueString,
		task.Description,
	)
	keyboard := PrepareSelectedTaskKeyboard()

	err = EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}
	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, keyboard)
	if err != nil {
		return err
	}

	if b.fsm.Current() == "updateTask" {
		err = b.fsm.Event(context.Background(), "selectedTask")
		if err != nil {
			return err
		}
	} else if b.fsm.Current() == "updateTodayTask" {
		err = b.fsm.Event(context.Background(), "selectedTodayTask")
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleDeleteTask(callBack *tgbotapi.CallbackQuery) error {
	accessToken, err := b.postgres.Get(callBack.Message.Chat.ID, databaseURL)
	if err != nil {
		return err
	}

	err, taskID := b.redis.Get(strconv.FormatInt(callBack.Message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.DeleteTask(taskID, accessToken)
	if err != nil {
		return err
	}

	taskList, err := b.client.GetAllTasks(accessToken)
	if err != nil {
		return err
	}

	tasksKeyboard := PrepareTaskListKeyboard(taskList)
	text := configs.Lexicon["response"]["task_list"]

	err = EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, tasksKeyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleDeleteTodayTask(callBack *tgbotapi.CallbackQuery) error {
	accessToken, err := b.postgres.Get(callBack.Message.Chat.ID, databaseURL)
	if err != nil {
		return err
	}

	err, taskID := b.redis.Get(strconv.FormatInt(callBack.Message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.DeleteTask(taskID, accessToken)
	if err != nil {
		return err
	}

	taskList, err := b.client.GetTodayTasks(accessToken)
	if err != nil {
		return err
	}

	tasksKeyboard := PrepareTaskListKeyboard(taskList)
	text := configs.Lexicon["response"]["task_list"]

	err = EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, tasksKeyboard)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "selectTodayTask")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTaskName(callBack *tgbotapi.CallbackQuery) error {
	text := configs.Lexicon["response"]["enter_new_task_name"]

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTaskName")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTaskPriority(callBack *tgbotapi.CallbackQuery) error {
	text := configs.Lexicon["response"]["enter_new_task_priority"]
	//keyboard := PrepareUpdatePriorityKeyboard

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	//err = EditMessageReplyMarkup(callBack.Message.Chat.ID, callBack.Message.MessageID, keyboard)

	err = b.fsm.Event(context.Background(), "updateTaskPriority")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTaskDueDate(callBack *tgbotapi.CallbackQuery) error {
	text := configs.Lexicon["response"]["enter_new_task_due_date"]

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTaskDueDate")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTaskDescription(callBack *tgbotapi.CallbackQuery) error {
	text := configs.Lexicon["response"]["enter_new_task_description"]

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTaskDescription")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTodayTaskName(callBack *tgbotapi.CallbackQuery) error {
	text := configs.Lexicon["response"]["enter_new_task_name"]

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTodayTaskName")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTodayTaskPriority(callBack *tgbotapi.CallbackQuery) error {
	text := configs.Lexicon["response"]["enter_new_task_priority"]

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTodayTaskPriority")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTodayTaskDueDate(callBack *tgbotapi.CallbackQuery) error {
	text := configs.Lexicon["response"]["enter_new_task_due_date"]

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTodayTaskDueDate")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdateTodayTaskDescription(callBack *tgbotapi.CallbackQuery) error {
	text := configs.Lexicon["response"]["enter_new_task_description"]

	err := EditMessageText(callBack.Message.Chat.ID, callBack.Message.MessageID, text)
	if err != nil {
		return err
	}

	err = b.fsm.Event(context.Background(), "updateTodayTaskDescription")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) updateTaskName(message *tgbotapi.Message) error {
	accessToken, err := b.postgres.Get(message.Chat.ID, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	err, taskID := b.redis.Get(strconv.FormatInt(message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.UpdateTaskName(taskID, accessToken, message.Text)
	if err != nil {
		return err
	}

	if b.fsm.Current() == "updateTaskName" {
		taskList, err := b.client.GetAllTasks(accessToken)
		if err != nil {
			return err
		}

		tasksKeyboard := PrepareTaskListKeyboard(taskList)
		text := configs.Lexicon["response"]["task_list"]

		err = SendMessage(message.Chat.ID, text, tasksKeyboard)
		if err != nil {
			return err
		}

		err = b.fsm.Event(context.Background(), "selectTask")
		if err != nil {
			return err
		}
	} else {
		taskList, err := b.client.GetTodayTasks(accessToken)
		if err != nil {
			return err
		}

		tasksKeyboard := PrepareTaskListKeyboard(taskList)
		text := configs.Lexicon["response"]["task_list"]

		err = SendMessage(message.Chat.ID, text, tasksKeyboard)
		if err != nil {
			return err
		}

		err = b.fsm.Event(context.Background(), "selectTodayTask")
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) updateTaskPriority(message *tgbotapi.Message) error {
	accessToken, err := b.postgres.Get(message.Chat.ID, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	err, taskID := b.redis.Get(strconv.FormatInt(message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.UpdateTaskPriority(taskID, accessToken, message.Text)

	if b.fsm.Current() == "updateTaskPriority" {
		taskList, err := b.client.GetAllTasks(accessToken)
		if err != nil {
			return err
		}

		tasksKeyboard := PrepareTaskListKeyboard(taskList)
		text := configs.Lexicon["response"]["task_list"]

		err = SendMessage(message.Chat.ID, text, tasksKeyboard)
		if err != nil {
			return err
		}

		err = b.fsm.Event(context.Background(), "selectTask")
		if err != nil {
			return err
		}
	} else {
		taskList, err := b.client.GetTodayTasks(accessToken)
		if err != nil {
			return err
		}

		tasksKeyboard := PrepareTaskListKeyboard(taskList)
		text := configs.Lexicon["response"]["task_list"]

		err = SendMessage(message.Chat.ID, text, tasksKeyboard)
		if err != nil {
			return err
		}

		err = b.fsm.Event(context.Background(), "selectTodayTask")
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) updateTaskDueDate(message *tgbotapi.Message) error {
	accessToken, err := b.postgres.Get(message.Chat.ID, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	err, taskID := b.redis.Get(strconv.FormatInt(message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.UpdateTaskDueDate(taskID, accessToken, message.Text)

	if b.fsm.Current() == "updateTaskDueDate" {
		taskList, err := b.client.GetAllTasks(accessToken)
		if err != nil {
			return err
		}

		tasksKeyboard := PrepareTaskListKeyboard(taskList)
		text := configs.Lexicon["response"]["task_list"]

		err = SendMessage(message.Chat.ID, text, tasksKeyboard)
		if err != nil {
			return err
		}

		err = b.fsm.Event(context.Background(), "selectTask")
		if err != nil {
			return err
		}
	} else {
		taskList, err := b.client.GetTodayTasks(accessToken)
		if err != nil {
			return err
		}

		tasksKeyboard := PrepareTaskListKeyboard(taskList)
		text := configs.Lexicon["response"]["task_list"]

		err = SendMessage(message.Chat.ID, text, tasksKeyboard)
		if err != nil {
			return err
		}

		err = b.fsm.Event(context.Background(), "selectTodayTask")
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) updateTaskDescription(message *tgbotapi.Message) error {
	accessToken, err := b.postgres.Get(message.Chat.ID, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	err, taskID := b.redis.Get(strconv.FormatInt(message.Chat.ID, 10))
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.UpdateTaskDescription(taskID, accessToken, message.Text)

	if b.fsm.Current() == "updateTaskDescription" {
		taskList, err := b.client.GetAllTasks(accessToken)
		if err != nil {
			return err
		}

		tasksKeyboard := PrepareTaskListKeyboard(taskList)
		text := configs.Lexicon["response"]["task_list"]

		err = SendMessage(message.Chat.ID, text, tasksKeyboard)
		if err != nil {
			return err
		}

		err = b.fsm.Event(context.Background(), "selectTask")
		if err != nil {
			return err
		}
	} else {
		taskList, err := b.client.GetTodayTasks(accessToken)
		if err != nil {
			return err
		}

		tasksKeyboard := PrepareTaskListKeyboard(taskList)
		text := configs.Lexicon["response"]["task_list"]

		err = SendMessage(message.Chat.ID, text, tasksKeyboard)
		if err != nil {
			return err
		}

		err = b.fsm.Event(context.Background(), "selectTodayTask")
		if err != nil {
			return err
		}
	}

	return nil
}
