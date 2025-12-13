package formatter

import (
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	formatter := NewFormatter("  ")
	if formatter == nil {
		t.Fatal("NewFormatter() returned nil")
	}
	if formatter.Indent != "  " {
		t.Errorf("NewFormatter() Indent = %q, want '  '", formatter.Indent)
	}
}

func TestFormatter_FormatList(t *testing.T) {
	tests := []struct {
		name   string
		indent string
		items  []string
		want   string
	}{
		{
			name:   "format list with items",
			indent: "",
			items:  []string{"item1", "item2", "item3"},
			want:   "1. item1\n2. item2\n3. item3\n",
		},
		{
			name:   "format list with indent",
			indent: "  ",
			items:  []string{"item1", "item2"},
			want:   "  1. item1\n  2. item2\n",
		},
		{
			name:   "format empty list",
			indent: "",
			items:  []string{},
			want:   "",
		},
		{
			name:   "format single item",
			indent: "",
			items:  []string{"single"},
			want:   "1. single\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter(tt.indent)
			got := f.FormatList(tt.items)
			if got != tt.want {
				t.Errorf("FormatList() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatter_FormatTable(t *testing.T) {
	f := NewFormatter("")

	t.Run("format table with data", func(t *testing.T) {
		headers := []string{"Name", "Age", "City"}
		rows := [][]string{
			{"Alice", "30", "New York"},
			{"Bob", "25", "London"},
		}

		// We can't easily test the output without capturing stdout,
		// but we can verify it doesn't panic
		f.FormatTable(headers, rows)
	})

	t.Run("format table with empty headers", func(t *testing.T) {
		headers := []string{}
		rows := [][]string{{"Alice", "30"}}

		// Should return early without panicking
		f.FormatTable(headers, rows)
	})

	t.Run("format table with mismatched row lengths", func(t *testing.T) {
		headers := []string{"Name", "Age"}
		rows := [][]string{
			{"Alice", "30", "Extra"},
			{"Bob"},
		}

		// Should handle gracefully
		f.FormatTable(headers, rows)
	})

	t.Run("format table with single row", func(t *testing.T) {
		headers := []string{"Name"}
		rows := [][]string{{"Alice"}}

		f.FormatTable(headers, rows)
	})
}

func TestFormatter_FormatOutput(t *testing.T) {
	f := NewFormatter("  ")

	// We can't easily test color output without capturing stdout,
	// but we can verify it doesn't panic
	t.Run("format stdout", func(t *testing.T) {
		f.FormatOutput("=>", "test output", true)
	})

	t.Run("format stderr", func(t *testing.T) {
		f.FormatOutput("!!", "test error", false)
	})
}

func TestFormatter_FormatHighlight(t *testing.T) {
	f := NewFormatter("")

	// Verify it doesn't panic
	f.FormatHighlight("Test Highlight")
}

func TestFormatter_FormatError(t *testing.T) {
	f := NewFormatter("  ")

	// Verify it doesn't panic
	f.FormatError("ERROR", "test error message")
}

func TestFormatter_FormatSuccess(t *testing.T) {
	f := NewFormatter("  ")

	// Verify it doesn't panic
	f.FormatSuccess("SUCCESS", "test success message")
}

func TestFormatter_FormatInfo(t *testing.T) {
	f := NewFormatter("  ")

	// Verify it doesn't panic
	f.FormatInfo("INFO", "test info message")
}

func TestFormatter_FormatWarning(t *testing.T) {
	f := NewFormatter("  ")

	// Verify it doesn't panic
	f.FormatWarning("WARN", "test warning message")
}

func TestFormatter_FormatDebug(t *testing.T) {
	f := NewFormatter("  ")

	// Verify it doesn't panic
	f.FormatDebug("DEBUG", "test debug message")
}

func TestFormatter_FormatTable_ColumnWidths(t *testing.T) {
	f := NewFormatter("")

	headers := []string{"Short", "Very Long Header Name", "Medium"}
	rows := [][]string{
		{"A", "B", "C"},
		{"Longer content here", "Short", "Even longer content in this cell"},
	}

	// This should calculate column widths correctly
	// We can't easily verify the exact output, but we can check it doesn't panic
	f.FormatTable(headers, rows)
}

func TestFormatter_FormatTable_EmptyRows(t *testing.T) {
	f := NewFormatter("")

	headers := []string{"Name", "Age"}
	rows := [][]string{}

	// Should handle empty rows gracefully
	f.FormatTable(headers, rows)
}

// Helper to check if string contains substring (for testing table output if needed)
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
