package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("  ")
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}
	if logger.formatter == nil {
		t.Error("NewLogger() formatter is nil")
	}
	if logger.level != LevelInfo {
		t.Errorf("NewLogger() level = %v, want %v", logger.level, LevelInfo)
	}
}

func TestLogger_SetLevel(t *testing.T) {
	logger := NewLogger("")

	tests := []struct {
		name  string
		level LogLevel
	}{
		{"set debug level", LevelDebug},
		{"set info level", LevelInfo},
		{"set warn level", LevelWarn},
		{"set error level", LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetLevel(tt.level)
			if logger.level != tt.level {
				t.Errorf("SetLevel() level = %v, want %v", logger.level, tt.level)
			}
		})
	}
}

func TestLogger_SetPrefix(t *testing.T) {
	logger := NewLogger("")

	tests := []struct {
		name   string
		prefix string
	}{
		{"set prefix", "TEST"},
		{"set empty prefix", ""},
		{"set long prefix", "Very Long Prefix Name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetPrefix(tt.prefix)
			if logger.prefix != tt.prefix {
				t.Errorf("SetPrefix() prefix = %q, want %q", logger.prefix, tt.prefix)
			}
		})
	}
}

func TestLogger_Debug(t *testing.T) {
	logger := NewLogger("")

	t.Run("debug with level debug", func(t *testing.T) {
		logger.SetLevel(LevelDebug)
		// Should not panic
		logger.Debug("test message")
		logger.Debug("formatted %s", "message")
	})

	t.Run("debug with level info", func(t *testing.T) {
		logger.SetLevel(LevelInfo)
		// Should not output (level > debug)
		logger.Debug("test message")
	})
}

func TestLogger_Info(t *testing.T) {
	logger := NewLogger("")

	// Should not panic
	logger.Info("test message")
	logger.Info("formatted %s", "message")
}

func TestLogger_Highlight(t *testing.T) {
	logger := NewLogger("")

	// Should not panic
	logger.Highlight("test highlight")
	logger.Highlight("formatted %s", "highlight")
}

func TestLogger_Warn(t *testing.T) {
	logger := NewLogger("")

	t.Run("warn with level warn", func(t *testing.T) {
		logger.SetLevel(LevelWarn)
		// Should not panic
		logger.Warn("test warning")
		logger.Warn("formatted %s", "warning")
	})

	t.Run("warn with level error", func(t *testing.T) {
		logger.SetLevel(LevelError)
		// Should not output (level > warn)
		logger.Warn("test warning")
	})
}

func TestLogger_Error(t *testing.T) {
	logger := NewLogger("")

	// Should not panic
	logger.Error("test error")
	logger.Error("formatted %s", "error")
}

func TestLogger_Success(t *testing.T) {
	logger := NewLogger("")

	// Should not panic
	logger.Success("test success")
	logger.Success("formatted %s", "success")
}

func TestLogger_Plain(t *testing.T) {
	logger := NewLogger("")

	// Should not panic
	logger.Plain("%s", "test plain")
	logger.Plain("formatted %s", "plain")
}

func TestLogger_Output(t *testing.T) {
	logger := NewLogger("")

	t.Run("output stdout", func(t *testing.T) {
		logger.Output("TAG", "test message", false)
	})

	t.Run("output stderr", func(t *testing.T) {
		logger.Output("TAG", "test error", true)
	})
}

func TestLogger_FormatTable(t *testing.T) {
	logger := NewLogger("")

	headers := []string{"Name", "Age"}
	rows := [][]string{
		{"Alice", "30"},
		{"Bob", "25"},
	}

	// Should not panic
	logger.FormatTable(headers, rows)
}

func TestLogger_LogLevels(t *testing.T) {
	// Test that log levels are properly ordered
	if LevelDebug >= LevelInfo {
		t.Error("LogLevel ordering: Debug should be < Info")
	}
	if LevelInfo >= LevelWarn {
		t.Error("LogLevel ordering: Info should be < Warn")
	}
	if LevelWarn >= LevelError {
		t.Error("LogLevel ordering: Warn should be < Error")
	}
}

func TestLogger_LogWithPrefix(t *testing.T) {
	logger := NewLogger("")
	logger.SetPrefix("TEST")

	// Should include prefix in output
	logger.Info("message")
	logger.Error("error message")
}

func TestLogger_LogLevelFiltering(t *testing.T) {
	logger := NewLogger("")

	// Set to error level - only errors should be logged
	logger.SetLevel(LevelError)

	// These should not output (but shouldn't panic)
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warning message")

	// This should output
	logger.Error("error message")
}

// Note: Testing Fatal() would cause the test to exit, so we skip it
// func TestLogger_Fatal(t *testing.T) {
//     logger := NewLogger("")
//     // This would exit the process, so we can't test it in unit tests
//     // logger.Fatal("fatal error")
// }
