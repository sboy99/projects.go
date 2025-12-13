package services

import (
	"strings"
	"testing"
)

func TestNewCLIService(t *testing.T) {
	service := NewCLIService()
	if service == nil {
		t.Fatal("NewCLIService() returned nil")
	}
	if service.logger == nil {
		t.Error("NewCLIService() logger is nil")
	}
}

func TestCLIService_IsValidCommand(t *testing.T) {
	service := NewCLIService()

	tests := []struct {
		name           string
		command        string
		isShellCommand bool
		wantErr        bool
	}{
		{"empty command", "", false, true},
		{"empty command with spaces", "   ", false, true},
		{"valid command", "echo", false, false},
		{"valid shell command", "echo hello", true, false},
		{"invalid executable", "nonexistentcommand12345", false, true},
		{"command with unbalanced single quotes", "'unclosed", true, true},
		{"command with unbalanced double quotes", "\"unclosed", true, true},
		{"command with balanced quotes", "'closed'", true, false},
		{"command with balanced double quotes", "\"closed\"", true, false},
		{"command with unbalanced parentheses", "(unclosed", true, true},
		{"command with balanced parentheses", "(closed)", true, false},
		{"command with quotes and parentheses", "echo 'test' (test)", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.IsValidCommand(tt.command, tt.isShellCommand)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidCommand(%q, %v) error = %v, wantErr %v", tt.command, tt.isShellCommand, err, tt.wantErr)
			}
		})
	}
}

func TestCLIService_containsShellOperators(t *testing.T) {
	service := NewCLIService()

	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{"pipe operator", "echo hello | grep test", true},
		{"and operator", "cmd1 && cmd2", true},
		{"or operator", "cmd1 || cmd2", true},
		{"semicolon", "cmd1; cmd2", true},
		{"redirect output", "cmd > file", true},
		{"append output", "cmd >> file", true},
		{"redirect input", "cmd < file", true},
		{"background", "cmd &", true},
		{"command substitution $()", "echo $(date)", true},
		{"command substitution backticks", "echo `date`", true},
		{"no operators", "echo hello", false},
		{"simple command", "ls -la", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.containsShellOperators(tt.command)
			if got != tt.want {
				t.Errorf("containsShellOperators(%q) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}

func TestCLIService_hasUnbalancedQuotes(t *testing.T) {
	service := NewCLIService()

	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{"unbalanced single quote", "'unclosed", true},
		{"unbalanced double quote", "\"unclosed", true},
		{"balanced single quotes", "'closed'", false},
		{"balanced double quotes", "\"closed\"", false},
		{"nested quotes", "'outer \"inner\"'", false},
		{"escaped quotes", "\\'test\\'", false},
		{"no quotes", "echo hello", false},
		{"mixed unbalanced", "'test\"", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.hasUnbalancedQuotes(tt.command)
			if got != tt.want {
				t.Errorf("hasUnbalancedQuotes(%q) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}

func TestCLIService_hasUnbalancedParentheses(t *testing.T) {
	service := NewCLIService()

	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{"unbalanced opening", "(unclosed", true},
		{"unbalanced closing", "unclosed)", true},
		{"balanced parentheses", "(closed)", false},
		{"nested balanced", "((nested))", false},
		{"in quotes", "'(test)'", false},
		{"no parentheses", "echo hello", false},
		{"multiple balanced", "(a) (b)", false},
		{"escaped parentheses", "\\(test\\)", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.hasUnbalancedParentheses(tt.command)
			if got != tt.want {
				t.Errorf("hasUnbalancedParentheses(%q) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}

func TestCLIService_Execute(t *testing.T) {
	service := NewCLIService()

	t.Run("execute simple command", func(t *testing.T) {
		result, err := service.Execute(ExecuteOptions{
			Command:   "echo hello",
			StreamLog: false,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result == nil {
			t.Fatal("Execute() returned nil result")
		}
		if result.ExitCode != 0 {
			t.Errorf("Execute() ExitCode = %d, want 0", result.ExitCode)
		}
		if !strings.Contains(result.Stdout, "hello") {
			t.Errorf("Execute() Stdout = %q, want to contain 'hello'", result.Stdout)
		}
	})

	t.Run("execute command with error", func(t *testing.T) {
		// Use a command that will fail
		result, err := service.Execute(ExecuteOptions{
			Command:   "false",
			StreamLog: false,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result == nil {
			t.Fatal("Execute() returned nil result")
		}
		if result.ExitCode == 0 {
			t.Error("Execute() ExitCode = 0, want non-zero")
		}
	})

	t.Run("execute empty command", func(t *testing.T) {
		_, err := service.Execute(ExecuteOptions{
			Command:   "",
			StreamLog: false,
		})
		if err == nil {
			t.Error("Execute() expected error for empty command, got nil")
		}
	})

	t.Run("execute with working directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		result, err := service.Execute(ExecuteOptions{
			Command:   "pwd",
			Dir:       tmpDir,
			StreamLog: false,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result == nil {
			t.Fatal("Execute() returned nil result")
		}
		if !strings.Contains(result.Stdout, tmpDir) {
			t.Errorf("Execute() Stdout = %q, want to contain %q", result.Stdout, tmpDir)
		}
	})

	t.Run("execute with environment variables", func(t *testing.T) {
		result, err := service.Execute(ExecuteOptions{
			Command:   "echo $TEST_VAR",
			Env:       []string{"TEST_VAR=test_value"},
			UseShell:  true,
			StreamLog: false,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result == nil {
			t.Fatal("Execute() returned nil result")
		}
		if !strings.Contains(result.Stdout, "test_value") {
			t.Errorf("Execute() Stdout = %q, want to contain 'test_value'", result.Stdout)
		}
	})

	t.Run("execute shell command with pipe", func(t *testing.T) {
		result, err := service.Execute(ExecuteOptions{
			Command:   "echo hello | tr '[:lower:]' '[:upper:]'",
			UseShell:  true,
			StreamLog: false,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result == nil {
			t.Fatal("Execute() returned nil result")
		}
		if !strings.Contains(strings.ToUpper(result.Stdout), "HELLO") {
			t.Errorf("Execute() Stdout = %q, want to contain 'HELLO'", result.Stdout)
		}
	})
}

func TestCLIService_Execute_StreamLog(t *testing.T) {
	service := NewCLIService()

	t.Run("execute with streaming", func(t *testing.T) {
		result, err := service.Execute(ExecuteOptions{
			Command:   "echo test",
			StreamLog: true,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result == nil {
			t.Fatal("Execute() returned nil result")
		}
		if result.ExitCode != 0 {
			t.Errorf("Execute() ExitCode = %d, want 0", result.ExitCode)
		}
	})
}

func TestExecuteResult(t *testing.T) {
	// Test that ExecuteResult struct is properly formed
	result := &ExecuteResult{
		ExitCode: 0,
		Stdout:   "test output",
		Stderr:   "test error",
		Error:    nil,
	}

	if result.ExitCode != 0 {
		t.Errorf("ExecuteResult.ExitCode = %d, want 0", result.ExitCode)
	}
	if result.Stdout != "test output" {
		t.Errorf("ExecuteResult.Stdout = %q, want 'test output'", result.Stdout)
	}
}

func TestExecuteOptions(t *testing.T) {
	opts := ExecuteOptions{
		Command:   "test",
		Dir:       "/tmp",
		Env:       []string{"KEY=value"},
		StreamLog: true,
		UseShell:  true,
		Shell:     "/bin/bash",
	}

	if opts.Command != "test" {
		t.Errorf("ExecuteOptions.Command = %q, want 'test'", opts.Command)
	}
	if opts.Dir != "/tmp" {
		t.Errorf("ExecuteOptions.Dir = %q, want '/tmp'", opts.Dir)
	}
	if len(opts.Env) != 1 {
		t.Errorf("ExecuteOptions.Env length = %d, want 1", len(opts.Env))
	}
	if !opts.StreamLog {
		t.Error("ExecuteOptions.StreamLog = false, want true")
	}
	if !opts.UseShell {
		t.Error("ExecuteOptions.UseShell = false, want true")
	}
}

// Test that commands with shell operators are automatically detected
func TestCLIService_Execute_AutoShellDetection(t *testing.T) {
	service := NewCLIService()

	t.Run("auto-detect pipe operator", func(t *testing.T) {
		result, err := service.Execute(ExecuteOptions{
			Command:   "echo hello | cat",
			StreamLog: false,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result == nil {
			t.Fatal("Execute() returned nil result")
		}
		// Should succeed because pipe is detected and shell is used
		if result.ExitCode != 0 {
			t.Errorf("Execute() ExitCode = %d, want 0", result.ExitCode)
		}
	})
}

// Test command that doesn't exist
func TestCLIService_Execute_NonExistentCommand(t *testing.T) {
	service := NewCLIService()

	// This should fail when not using shell
	_, err := service.Execute(ExecuteOptions{
		Command:   "nonexistentcommand12345xyz",
		StreamLog: false,
		UseShell:  false,
	})
	// Error is expected, but it might be in result.Error instead
	// Let's check both
	if err == nil {
		// Check if it's in the result
		result, _ := service.Execute(ExecuteOptions{
			Command:   "nonexistentcommand12345xyz",
			StreamLog: false,
			UseShell:  false,
		})
		if result != nil && result.Error == nil && result.ExitCode == 0 {
			t.Error("Execute() should fail for non-existent command")
		}
	}
}
