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

func TransformToTableData[T any](data map[string]T) (headers []string, rows [][]string) {
	headers = make([]string, 0, len(data))
	rows = make([][]string, 0, len(data))
	for key, value := range data {
		headers = append(headers, key)
		rows = append(rows, []string{fmt.Sprintf("%v", value)})
	}
	return headers, rows
}
