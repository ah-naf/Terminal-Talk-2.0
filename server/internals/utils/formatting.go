package utils

import "fmt"

// FormatJoinMessage formats the join message
func FormatJoinMessage(username string) string {
	return fmt.Sprintf("\033[1;36m%s has entered the chat! 🎉\033[0m", username)
}

// FormatLeaveMessage formats the leave message
func FormatLeaveMessage(username string) string {
	return fmt.Sprintf("\033[1;31m%s left the chat!\033[0m", username)
}

// FormatChatMessage formats the chat message
func FormatChatMessage(username, message string) string {
	return fmt.Sprintf("\033[1;32m%s > %s\033[0m", username, message)
}
