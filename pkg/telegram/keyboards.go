package telegram

import (
	configs "Todoist-bot/pkg/config"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	todoist "github.com/hhhanafubukiii/go-todoist-sdk"
)

var (
	MainKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["today_tasks"]),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["new_task"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["all_tasks"]),
		),
	)
	CancelKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cancel ❌"),
		),
	)
	CancelSkipKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cancel ❌"),
			tgbotapi.NewKeyboardButton("Skip ➡️"),
		),
	)
	PriorityKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["priority_1"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["priority_2"]),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["priority_3"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["priority_4"]),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cancel ❌"),
			tgbotapi.NewKeyboardButton("Skip ➡️"),
		),
	)
	DueDateKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["today"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["tomorrow"]),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["next_week"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["next_month"]),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cancel ❌"),
			tgbotapi.NewKeyboardButton("Skip ➡️"),
		),
	)

	UpdatePriorityKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["priority_1"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["priority_2"]),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["priority_3"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["priority_4"]),
		),
	)

	UpdateDueDateKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["today"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["tomorrow"]),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["next_week"]),
			tgbotapi.NewKeyboardButton(configs.Lexicon["buttons"]["next_month"]),
		),
	)

	cancelButtonCallBack = configs.Lexicon["buttons"]["cancel"]
	backButtonCallBack   = configs.Lexicon["buttons"]["back"]
)

func PrepareTaskListKeyboard(taskList []todoist.Task) string {
	taskListInlineKeyboard := make(map[string][][]map[string]string, 1)
	taskListInlineKeyboard["inline_keyboard"] = make([][]map[string]string, len(taskList)+1)

	for k, task := range taskList {
		taskListInlineKeyboard["inline_keyboard"][k] = []map[string]string{{"text": task.Name, "callback_data": task.Id}}
	}

	taskListInlineKeyboard["inline_keyboard"][len(taskList)] = append([]map[string]string{{"text": configs.Lexicon["buttons"]["cancel"],
		"callback_data": cancelButtonCallBack}})

	marshalTaskListInlineKeyboard, _ := json.Marshal(taskListInlineKeyboard)

	return string(marshalTaskListInlineKeyboard)
}

func PrepareSelectedTaskKeyboard() string {
	completeButtonCallBack := configs.Lexicon["buttons"]["complete_task"]
	updateButtonCallBack := configs.Lexicon["buttons"]["update_task"]
	deleteButtonCallBack := configs.Lexicon["buttons"]["delete_task"]

	taskKeyboard := fmt.Sprintf(`{"inline_keyboard":[[{"text": "%s", "callback_data": "%s"}], [{"text": "%s", "callback_data": "%s"}, {"text": "%s", "callback_data": "%s"}], [{"text": "⬅️", "callback_data": "%s"}]]}`,
		configs.Lexicon["buttons"]["complete_task"],
		completeButtonCallBack,
		configs.Lexicon["buttons"]["update_task"],
		updateButtonCallBack,
		configs.Lexicon["buttons"]["delete_task"],
		deleteButtonCallBack,
		backButtonCallBack,
	)

	return taskKeyboard
}

func PrepareUpdateTaskKeyboard() string {
	nameButtonCallBack := configs.Lexicon["buttons"]["name"]
	priorityButtonCallBack := configs.Lexicon["buttons"]["priority"]
	dueDateButtonCallBack := configs.Lexicon["buttons"]["due_date"]
	descriptionButtonCallBack := configs.Lexicon["buttons"]["description"]

	updateTaskKeyboard := fmt.Sprintf(`{"inline_keyboard": [[{"text": "Name","callback_data": "%s"},{"text": "Priority","callback_data": "%s"}],[{"text": "Due date","callback_data": "%s"},{"text": "Description","callback_data": "%s"}], [{"text": "⬅️", "callback_data": "%s"}]]}`,
		nameButtonCallBack,
		priorityButtonCallBack,
		dueDateButtonCallBack,
		descriptionButtonCallBack,
		backButtonCallBack,
	)

	return updateTaskKeyboard
}

//func PrepareUpdatePriorityKeyboard() string {
//
//}
//
//func PrepareUpdateDueDateKeyboard() string {
//
//}
