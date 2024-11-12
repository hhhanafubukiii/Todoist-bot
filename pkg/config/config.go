package configs

var Lexicon = map[string]map[string]string{
	"response": {
		"start":              "Hello! In order for me to interact with your Todoist account, you need to give me access to it. For this link follow the link:\n%s",
		"already_authorized": "You're already signed in with your Todoist account.",
		"unknown_command":    "I dont now such a command.",
	},
}

var databaseRequests = map[string]string{
	"save": "INSERT INTO tokens(chat_id, access_token) VALUES (%d, '%s')",
	"get":  "SELECT access_token FROM tokens WHERE chat_id = %d",
}
