package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/sboy99/projects.go/tusk/pkg/logger"
)

// CLIService handles command execution with streaming output
type CLIService struct {
	logger *logger.Logger
}

// NewCLIService creates a new CLI service instance
func NewCLIService() *CLIService {
	return &CLIService{
		logger: logger.NewLogger("  "),
	}
}

// ExecuteOptions contains options for command execution
type ExecuteOptions struct {
	Command   string
	Dir       string
	Env       []string
	StreamLog bool   // Whether to stream logs in real-time
	UseShell  bool   // Whether to execute via shell (supports pipes, redirects, etc.)
	Shell     string // Shell to use (default: /bin/sh on Unix)
}

// ExecuteResult contains the result of command execution
type ExecuteResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
}

// Execute executes a command and optionally streams the output
func (c *CLIService) Execute(opts ExecuteOptions) (*ExecuteResult, error) {
	hasShellOps := c.containsShellOperators(opts.Command)
	isShellCmd := opts.UseShell || hasShellOps

	var cmd *exec.Cmd

	// Execute via shell if requested or if command contains shell operators
	if isShellCmd {
		shell := opts.Shell
		if shell == "" {
			shell = "/bin/sh"
		}
		cmd = exec.Command(shell, "-c", opts.Command)
	} else {
		// Parse command into executable and arguments
		parts := strings.Fields(opts.Command)
		if len(parts) == 0 {
			return nil, fmt.Errorf("empty command")
		}
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	// Set working directory if provided
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}

	// Set environment variables if provided
	if len(opts.Env) > 0 {
		cmd.Env = append(os.Environ(), opts.Env...)
	}

	var stdoutBuf, stderrBuf strings.Builder
	var stdoutPipe, stderrPipe io.ReadCloser
	var err error
	var wg sync.WaitGroup

	// Setup stdout pipe
	stdoutPipe, err = cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Setup stderr pipe
	stderrPipe, err = cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Stream output if requested
	if opts.StreamLog {
		// Stream stdout
		wg.Go(func() {
			c.streamOutput(stdoutPipe, &stdoutBuf, "=>", true)
		})
		// Stream stderr
		wg.Go(func() {
			c.streamOutput(stderrPipe, &stderrBuf, "!!", false)
		})
	} else {
		// Collect output without streaming
		wg.Go(func() {
			io.Copy(&stdoutBuf, stdoutPipe)
		})
		wg.Go(func() {
			io.Copy(&stderrBuf, stderrPipe)
		})
	}

	// Wait for command to complete
	err = cmd.Wait()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return &ExecuteResult{
				ExitCode: -1,
				Stdout:   stdoutBuf.String(),
				Stderr:   stderrBuf.String(),
				Error:    err,
			}, err
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return &ExecuteResult{
		ExitCode: exitCode,
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		Error:    err,
	}, nil
}

// isValidCommand validates a command before execution
func (c *CLIService) IsValidCommand(command string, isShellCommand bool) error {
	// Check if command is empty
	command = strings.TrimSpace(command)
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// For shell commands, basic syntax validation
	if isShellCommand {
		// Check for unbalanced quotes
		if c.hasUnbalancedQuotes(command) {
			return fmt.Errorf("command has unbalanced quotes")
		}
		// Check for unbalanced parentheses
		if c.hasUnbalancedParentheses(command) {
			return fmt.Errorf("command has unbalanced parentheses")
		}
		return nil
	}

	// For non-shell commands, check if executable exists
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("command cannot be empty")
	}

	executable := parts[0]
	if _, err := exec.LookPath(executable); err != nil {
		return err
	}

	return nil
}

// streamOutput streams output from a reader with formatted output
func (c *CLIService) streamOutput(reader io.ReadCloser, buf *strings.Builder, prefix string, isStdout bool) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		buf.WriteString(line + "\n")

		// Format and print the line using logger
		c.logger.Output(prefix, line, !isStdout)
	}
}

// hasUnbalancedQuotes checks if a command has unbalanced quotes
func (c *CLIService) hasUnbalancedQuotes(command string) bool {
	singleQuoteCount := 0
	doubleQuoteCount := 0
	escaped := false

	for _, char := range command {
		if escaped {
			escaped = false
			continue
		}
		if char == '\\' {
			escaped = true
			continue
		}
		if char == '\'' && doubleQuoteCount%2 == 0 {
			singleQuoteCount++
		}
		if char == '"' && singleQuoteCount%2 == 0 {
			doubleQuoteCount++
		}
	}

	return singleQuoteCount%2 != 0 || doubleQuoteCount%2 != 0
}

// hasUnbalancedParentheses checks if a command has unbalanced parentheses
func (c *CLIService) hasUnbalancedParentheses(command string) bool {
	parenCount := 0
	inString := false
	stringChar := byte(0)
	escaped := false

	for i := 0; i < len(command); i++ {
		char := command[i]
		if escaped {
			escaped = false
			continue
		}
		if char == '\\' {
			escaped = true
			continue
		}
		if (char == '\'' || char == '"') && !inString {
			inString = true
			stringChar = char
		} else if char == stringChar && inString {
			inString = false
			stringChar = 0
		}
		if !inString {
			if char == '(' {
				parenCount++
			} else if char == ')' {
				parenCount--
				if parenCount < 0 {
					return true
				}
			}
		}
	}

	return parenCount != 0
}

// containsShellOperators checks if command contains shell operators
func (c *CLIService) containsShellOperators(command string) bool {
	shellOps := []string{"|", "&&", "||", ";", ">", ">>", "<", "&", "$(", "`"}
	for _, op := range shellOps {
		if strings.Contains(command, op) {
			return true
		}
	}
	return false
}
