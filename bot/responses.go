package bot

// BuildCommandAcceptedResponse returns стандартный ответ для BotX.
func BuildCommandAcceptedResponse() map[string]string {
	return map[string]string{"status": "accepted"}
}

// BuildBotDisabledResponse returns response for disabled bot.
func BuildBotDisabledResponse() map[string]string {
	return map[string]string{"status": "bot_disabled"}
}
