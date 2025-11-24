package formatter

import (
	"fmt"
	"strings"
)

// Formatter provides formatting utilities
type Formatter struct {
	Indent string
}

// NewFormatter creates a new formatter instance
func NewFormatter(indent string) *Formatter {
	return &Formatter{
		Indent: indent,
	}
}

// FormatList formats a list of items
func (f *Formatter) FormatList(items []string) string {
	var sb strings.Builder
	for i, item := range items {
		sb.WriteString(fmt.Sprintf("%s%d. %s\n", f.Indent, i+1, item))
	}
	return sb.String()
}

// FormatTable formats data as a table
func (f *Formatter) FormatTable(headers []string, rows [][]string) string {
	var sb strings.Builder
	
	// Write headers
	sb.WriteString(strings.Join(headers, "\t"))
	sb.WriteString("\n")
	
	// Write separator
	sb.WriteString(strings.Repeat("-", len(strings.Join(headers, "\t"))))
	sb.WriteString("\n")
	
	// Write rows
	for _, row := range rows {
		sb.WriteString(strings.Join(row, "\t"))
		sb.WriteString("\n")
	}
	
	return sb.String()
}

