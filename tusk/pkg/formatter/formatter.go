package formatter

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
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

// FormatTable formats data as a table (Docker-style)
func (f *Formatter) FormatTable(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	// Find maximum width for each column
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				if len(cell) > colWidths[i] {
					colWidths[i] = len(cell)
				}
			}
		}
	}

	// Add padding (minimum 2 spaces between columns)
	for i := range colWidths {
		colWidths[i] += 2
	}

	var sb strings.Builder

	// Write headers
	for i, header := range headers {
		if i > 0 {
			sb.WriteString("  ") // 2 spaces between columns
		}
		sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], header))
	}
	sb.WriteString("\n")

	// Write separator (Docker-style: spaces and dashes)
	for i, width := range colWidths {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(strings.Repeat("-", width))
	}
	sb.WriteString("\n")

	// Write rows
	for _, row := range rows {
		for i, cell := range row {
			if i >= len(colWidths) {
				break
			}
			if i > 0 {
				sb.WriteString("  ") // 2 spaces between columns
			}
			sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], cell))
		}
		sb.WriteString("\n")
	}

	fmt.Print(sb.String())
}

// FormatOutput formats CLI output with prefix and color coding
func (f *Formatter) FormatOutput(prefix, line string, isStdout bool) {
	if isStdout {
		// Format stdout with prefix
		color.New(color.FgGreen).Printf("%s[%s] %s\n", f.Indent, prefix, line)
	} else {
		// Format stderr with prefix (could use different formatting)
		color.New(color.FgRed).Fprintf(os.Stderr, "%s[%s] %s\n", f.Indent, prefix, line)
	}
}

func (f *Formatter) FormatHighlight(line string) {
	color.New(color.FgMagenta).Printf("[+] %s\n", line)
}

// FormatError formats an error message with prefix
func (f *Formatter) FormatError(prefix, line string) {
	// Format stderr with prefix (could use different formatting)
	color.New(color.FgRed).Fprintf(os.Stderr, "%s[%s] %s\n", f.Indent, prefix, line)
}

// FormatSuccess formats a success message with prefix
func (f *Formatter) FormatSuccess(prefix, line string) {
	color.New(color.FgGreen).Printf("%s[%s] %s\n", f.Indent, prefix, line)
}

func (f *Formatter) FormatInfo(prefix, line string) {
	color.New(color.FgBlue).Printf("%s[%s] %s\n", f.Indent, prefix, line)
}

func (f *Formatter) FormatWarning(prefix, line string) {
	color.New(color.FgYellow).Printf("%s[%s] %s\n", f.Indent, prefix, line)
}

func (f *Formatter) FormatDebug(prefix, line string) {
	color.New(color.FgBlack).Printf("%s[%s] %s\n", f.Indent, prefix, line)
}

func (f *Formatter) FormatPlain(prefix, line string) {
	color.New(color.FgWhite).Printf("%s[%s] %s\n", f.Indent, prefix, line)
}
