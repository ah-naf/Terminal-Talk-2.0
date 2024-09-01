package utils

import "fmt"

// ANSI escape codes for text formatting
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	FgCyan      = "\033[1;36m" // Bold Cyan
	FgRed       = "\033[1;31m" // Bold Red
	FgGreen     = "\033[1;32m" // Bold Green
	FgYellow    = "\033[1;33m" // Bold Yellow
	Italic      = "\033[3m" 
)

// FormatJoinMessage formats the join message
func FormatJoinMessage(username string) string {
	return fmt.Sprintf("%s%s has entered the chat! 🎉%s", FgCyan, username, Reset)
}

// FormatLeaveMessage formats the leave message
func FormatLeaveMessage(username string) string {
	return fmt.Sprintf("%s%s left the chat!%s", FgRed, username, Reset)
}

// FormatChatMessage formats the chat message
func FormatChatMessage(username, message string) string {
	return fmt.Sprintf("%s%s > %s%s", FgGreen, username, message, Reset)
}

// FormatSuccessMessage formats a success message
func FormatSuccessMessage(message string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, FgGreen, message, Reset)
}

// FormatErrorMessage formats an error message
func FormatErrorMessage(message string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, FgRed, message, Reset)
}

// FormatWarningMessage formats a warning message
func FormatWarningMessage(message string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, FgYellow, message, Reset)
}
