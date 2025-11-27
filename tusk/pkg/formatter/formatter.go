package formatter

import (
	"fmt"
	"os"
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

// FormatOutput formats CLI output with prefix and color coding
func (f *Formatter) FormatOutput(prefix, line string, isStdout bool) {
	if isStdout {
		// Format stdout with prefix
		fmt.Printf("%s[%s] %s\n", f.Indent, prefix, line)
	} else {
		// Format stderr with prefix (could use different formatting)
		fmt.Fprintf(os.Stderr, "%s[%s] %s\n", f.Indent, prefix, line)
	}
}
