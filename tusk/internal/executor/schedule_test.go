package executor

import (
	"testing"
)

func TestNewScheduleExecutor(t *testing.T) {
	executor := NewScheduleExecutor()
	if executor == nil {
		t.Fatal("NewScheduleExecutor() returned nil")
	}
	if executor.scheduleService == nil {
		t.Error("NewScheduleExecutor() scheduleService is nil")
	}
	if executor.sudoService == nil {
		t.Error("NewScheduleExecutor() sudoService is nil")
	}
	if executor.logger == nil {
		t.Error("NewScheduleExecutor() logger is nil")
	}
}

func TestScheduleExecutor_Create(t *testing.T) {
	executor := NewScheduleExecutor()

	// Note: Full testing of Create() would require:
	// 1. Mocking SudoService to avoid actual sudo prompts
	// 2. Mocking ScheduleService to avoid actual systemd operations
	// 3. Verifying error handling

	// For now, we test that the method exists and doesn't panic
	// In a real scenario, we'd use dependency injection and mocks
	t.Run("create method exists", func(t *testing.T) {
		// This would normally require sudo and systemd
		// So we just verify the method is callable
		// executor.Create("test", "echo hello", "1m")
		_ = executor.Create // Verify method exists
	})
}

func TestScheduleExecutor_List(t *testing.T) {
	executor := NewScheduleExecutor()

	// Similar to Create, full testing would require mocking
	t.Run("list method exists", func(t *testing.T) {
		_ = executor.List // Verify method exists
	})
}

func TestScheduleExecutor_Delete(t *testing.T) {
	executor := NewScheduleExecutor()

	// Similar to Create, full testing would require mocking
	t.Run("delete method exists", func(t *testing.T) {
		_ = executor.Delete // Verify method exists
	})
}

// Note: To properly test these methods with unit tests, we would need:
// 1. Dependency injection in ScheduleExecutor
// 2. Interface-based design for ScheduleService and SudoService
// 3. Mock implementations for testing
//
// Example of how it could be structured:
//
// type ScheduleServiceInterface interface {
//     Create(name, command, interval string) error
//     List() error
//     Delete(name string) error
// }
//
// type ScheduleExecutor struct {
//     scheduleService ScheduleServiceInterface
//     sudoService     SudoServiceInterface
//     logger          *logger.Logger
// }
//
// Then in tests:
// executor := &ScheduleExecutor{
//     scheduleService: &MockScheduleService{},
//     sudoService:     &MockSudoService{},
//     logger:          logger.NewLogger(""),
// }
