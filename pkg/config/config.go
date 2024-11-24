package configs

var Lexicon = map[string]map[string]string{
	"response": {
		"start":              "Hi! I need to access your Todoist account to get started. Please click here to sign in:\n%s",
		"already_authorized": "You are already authorized!",
		"unknown_command":    "Please choose a command from the command list.",
		"new_task":           "OK, a new task! Send me a message in this format:\ntask_name;priority;deadline;description",
		"add_task":           "Success! Task added to your list",
	},
}

var DatabaseRequests = map[string]string{
	"save": "INSERT INTO tokens(chat_id, access_token) VALUES (%d, '%s')",
	"get":  "SELECT access_token FROM tokens WHERE chat_id = %d",
}
