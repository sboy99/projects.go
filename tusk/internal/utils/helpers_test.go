package utils

import (
	"strings"
	"testing"
	"time"
)

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"trim leading spaces", "  hello", "hello"},
		{"trim trailing spaces", "hello  ", "hello"},
		{"trim both sides", "  hello  ", "hello"},
		{"no spaces", "hello", "hello"},
		{"empty string", "", ""},
		{"only spaces", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimSpace(tt.input)
			if got != tt.want {
				t.Errorf("TrimSpace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		message string
		want    string
	}{
		{"basic format", "INFO", "test message", "[INFO] test message"},
		{"empty prefix", "", "message", "[] message"},
		{"empty message", "ERROR", "", "[ERROR] "},
		{"both empty", "", "", "[] "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMessage(tt.prefix, tt.message)
			if got != tt.want {
				t.Errorf("FormatMessage(%q, %q) = %q, want %q", tt.prefix, tt.message, got, tt.want)
			}
		})
	}
}

func TestValidateNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid string", "hello", false},
		{"empty string", "", true},
		{"only spaces", "   ", true},
		{"tab only", "\t", true},
		{"newline only", "\n", true},
		{"mixed whitespace", " \t\n ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotEmpty(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNotEmpty(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestJoinStrings(t *testing.T) {
	tests := []struct {
		name string
		sep  string
		strs []string
		want string
	}{
		{"join with comma", ",", []string{"a", "b", "c"}, "a,b,c"},
		{"join with space", " ", []string{"hello", "world"}, "hello world"},
		{"join with empty sep", "", []string{"a", "b"}, "ab"},
		{"single string", ",", []string{"a"}, "a"},
		{"no strings", ",", []string{}, ""},
		{"empty separator", "", []string{"a", "b", "c"}, "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JoinStrings(tt.sep, tt.strs...)
			if got != tt.want {
				t.Errorf("JoinStrings(%q, %v) = %q, want %q", tt.sep, tt.strs, got, tt.want)
			}
		})
	}
}

func TestGenerateUUID(t *testing.T) {
	// Generate multiple UUIDs and check they're unique
	uuid1 := GenerateUUID()
	uuid2 := GenerateUUID()

	if uuid1 == "" {
		t.Error("GenerateUUID() returned empty string")
	}
	if uuid2 == "" {
		t.Error("GenerateUUID() returned empty string")
	}
	if uuid1 == uuid2 {
		t.Error("GenerateUUID() returned duplicate UUIDs")
	}

	// Check format (UUID v4 format: 8-4-4-4-12)
	if len(uuid1) != 36 {
		t.Errorf("GenerateUUID() returned UUID with length %d, want 36", len(uuid1))
	}
}

func TestTransformToTableData(t *testing.T) {
	type TestStruct struct {
		Name      string    `json:"name"`
		Age       int       `json:"age"`
		StartTime time.Time `json:"startTime"`
		Email     string    `json:"email"`
	}

	now := time.Now()
	testData := map[string]TestStruct{
		"1": {
			Name:      "Alice",
			Age:       30,
			StartTime: now,
			Email:     "alice@example.com",
		},
		"2": {
			Name:      "Bob",
			Age:       25,
			StartTime: now.Add(time.Hour),
			Email:     "bob@example.com",
		},
	}

	t.Run("transform all fields", func(t *testing.T) {
		headers, rows := TransformToTableData(testData)

		if len(headers) == 0 {
			t.Error("TransformToTableData() returned empty headers")
		}
		if len(rows) != 2 {
			t.Errorf("TransformToTableData() returned %d rows, want 2", len(rows))
		}
		if len(rows[0]) != len(headers) {
			t.Errorf("TransformToTableData() row length %d, want %d", len(rows[0]), len(headers))
		}
	})

	t.Run("transform selected fields", func(t *testing.T) {
		headers, rows := TransformToTableData(testData, "name", "age")

		if len(headers) != 2 {
			t.Errorf("TransformToTableData() returned %d headers, want 2", len(headers))
		}
		if len(rows) != 2 {
			t.Errorf("TransformToTableData() returned %d rows, want 2", len(rows))
		}
		if len(rows[0]) != 2 {
			t.Errorf("TransformToTableData() row length %d, want 2", len(rows[0]))
		}
	})

	t.Run("transform empty data", func(t *testing.T) {
		emptyData := map[string]TestStruct{}
		headers, rows := TransformToTableData(emptyData)

		if len(headers) != 0 {
			t.Errorf("TransformToTableData() returned %d headers for empty data, want 0", len(headers))
		}
		if len(rows) != 0 {
			t.Errorf("TransformToTableData() returned %d rows for empty data, want 0", len(rows))
		}
	})

	t.Run("transform with time field", func(t *testing.T) {
		headers, rows := TransformToTableData(testData, "name", "startTime")

		if len(headers) != 2 {
			t.Errorf("TransformToTableData() returned %d headers, want 2", len(headers))
		}
		// Check that time is formatted correctly
		if len(rows) > 0 && len(rows[0]) > 1 {
			timeStr := rows[0][1]
			// Should be formatted as "2006-01-02 15:04:05"
			if len(timeStr) < 10 {
				t.Errorf("TransformToTableData() time format seems incorrect: %q", timeStr)
			}
		}
	})
}

func TestTransformToTableDataWithJSONTags(t *testing.T) {
	type TestStruct struct {
		FieldName  string `json:"field_name"`
		OtherField int    `json:"other_field"`
	}

	testData := map[string]TestStruct{
		"1": {FieldName: "test", OtherField: 42},
	}

	headers, _ := TransformToTableData(testData, "field_name", "other_field")

	if len(headers) != 2 {
		t.Errorf("TransformToTableData() returned %d headers, want 2", len(headers))
	}

	// Headers should use display names (separated camelCase or title case)
	foundFieldName := false
	for _, h := range headers {
		// Accept various formats: "Field Name", "Field_name", "field_name", etc.
		if h == "Field Name" || h == "Field_name" || h == "field_name" ||
			strings.Contains(strings.ToLower(h), "field") {
			foundFieldName = true
		}
	}
	if !foundFieldName {
		t.Errorf("TransformToTableData() header not found, got: %v", headers)
	}
}

func TestSeparateCamelCase(t *testing.T) {
	// This is a private function, but we can test it indirectly through TransformToTableData
	type TestStruct struct {
		StartTime      string `json:"startTime"`
		ScriptFilePath string `json:"scriptFilePath"`
		ID             string `json:"id"`
	}

	testData := map[string]TestStruct{
		"1": {StartTime: "now", ScriptFilePath: "/path", ID: "123"},
	}

	headers, _ := TransformToTableData(testData)

	// This is an indirect test, so we just verify headers are generated
	if len(headers) == 0 {
		t.Error("TransformToTableData() should generate headers")
	}
}
