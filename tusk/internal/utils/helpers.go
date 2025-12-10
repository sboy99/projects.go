package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// separateCamelCase separates camelCase words by inserting spaces
// Example: "StartTime" -> "Start Time", "ScriptFilePath" -> "Script File Path"
func separateCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder
	result.Grow(len(s) + 5) // Pre-allocate with some extra space

	for i, r := range s {
		// Insert space before uppercase letters (except the first character)
		if i > 0 && r >= 'A' && r <= 'Z' {
			// Check if previous character was also uppercase (to handle acronyms like "ID")
			prevRune := rune(s[i-1])
			if !(prevRune >= 'A' && prevRune <= 'Z') {
				result.WriteRune(' ')
			}
		}
		result.WriteRune(r)
	}

	return result.String()
}

func TransformToTableData[T any](data map[string]T, selectedFields ...string) (headers []string, rows [][]string) {
	if len(data) == 0 {
		return []string{}, [][]string{}
	}

	// Get the type of T to extract field names
	var zero T
	valueType := reflect.TypeOf(zero)

	// Build a map of field names (json tag or struct name) to field indices
	fieldMap := make(map[string]int)
	fieldDisplayNames := make(map[int]string)

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		// Get json tag name
		jsonTag := field.Tag.Get("json")
		var fieldKey string
		if jsonTag != "" && jsonTag != "-" {
			// Remove json tag options (e.g., "id,omitempty" -> "id")
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				fieldKey = jsonTag[:idx]
			} else {
				fieldKey = jsonTag
			}
		} else {
			fieldKey = field.Name
		}

		// Map both json tag name and struct field name to the same index
		fieldMap[fieldKey] = i
		fieldMap[field.Name] = i

		// Store display name (prefer json tag, fallback to field name)
		displayName := fieldKey
		// Separate camelCase words (e.g., "StartTime" -> "Start Time")
		displayName = separateCamelCase(displayName)
		caser := cases.Title(language.English)
		fieldDisplayNames[i] = caser.String(displayName)
	}

	// Determine which fields to include
	var fieldIndices []int
	if len(selectedFields) > 0 {
		// Use selected fields
		for _, selectedField := range selectedFields {
			if idx, ok := fieldMap[selectedField]; ok {
				fieldIndices = append(fieldIndices, idx)
			}
		}
	} else {
		// Use all fields (default)
		for i := 0; i < valueType.NumField(); i++ {
			fieldIndices = append(fieldIndices, i)
		}
	}

	// Extract headers from selected fields
	headers = make([]string, 0, len(fieldIndices))
	for _, idx := range fieldIndices {
		headers = append(headers, fieldDisplayNames[idx])
	}

	// Extract rows from struct values
	rows = make([][]string, 0, len(data))
	for _, value := range data {
		row := make([]string, 0, len(fieldIndices))
		valueReflect := reflect.ValueOf(value)

		for _, idx := range fieldIndices {
			fieldValue := valueReflect.Field(idx)
			var cellValue string

			// Handle different types appropriately
			switch fieldValue.Kind() {
			case reflect.String:
				cellValue = fieldValue.String()
			case reflect.Struct:
				// Special handling for time.Time
				if t, ok := fieldValue.Interface().(time.Time); ok {
					cellValue = t.Format("2006-01-02 15:04:05")
				} else {
					cellValue = fmt.Sprintf("%v", fieldValue.Interface())
				}
			default:
				cellValue = fmt.Sprintf("%v", fieldValue.Interface())
			}

			row = append(row, cellValue)
		}
		rows = append(rows, row)
	}

	return headers, rows
}
