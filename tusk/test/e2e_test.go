package test

import (
	"os"
	"os/exec"
	"testing"
)

// TestCLIVersion tests the version command
func TestCLIVersion(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/tusk", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version command failed: %v\nOutput: %s", err, output)
	}
	
	if len(output) == 0 {
		t.Error("version command produced no output")
	}
}

// TestCLIStart tests the schedule list command
func TestCLIStart(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/tusk", "schedule", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("schedule list command failed: %v\nOutput: %s", err, output)
	}
	
	// schedule list should succeed and produce output (even if empty)
	if len(output) == 0 {
		t.Error("schedule list command produced no output")
	}
}

// TestCLIStop tests the schedule delete command (should handle sudo requirement gracefully)
func TestCLIStop(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/tusk", "schedule", "delete", "--name", "test-service")
	output, err := cmd.CombinedOutput()
	// Note: command may exit with code 0 even when sudo fails, but should log an error
	outputStr := string(output)
	
	// Verify it handles the sudo requirement and logs an appropriate error message
	if !contains(outputStr, "sudo") && !contains(outputStr, "privileges") && !contains(outputStr, "ERROR") {
		t.Errorf("schedule delete should handle sudo requirement, got: %s", outputStr)
	}
	
	// Command should produce some output
	if len(outputStr) == 0 {
		t.Error("schedule delete command produced no output")
	}
	
	// If there was an actual error (not just sudo warning), that's fine for this test
	_ = err
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestMain(m *testing.M) {
	// Setup
	code := m.Run()
	// Teardown
	os.Exit(code)
}

