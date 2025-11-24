package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// TrimSpace trims whitespace from a string
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// FormatMessage formats a message with a prefix
func FormatMessage(prefix, message string) string {
	return fmt.Sprintf("[%s] %s", prefix, message)
}

// ValidateNotEmpty checks if a string is not empty
func ValidateNotEmpty(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("string cannot be empty")
	}
	return nil
}

// JoinStrings joins multiple strings with a separator
func JoinStrings(sep string, strs ...string) string {
	return strings.Join(strs, sep)
}

// GenerateID generates a random ID
func GenerateUUID() string {
	return uuid.New().String()
}
