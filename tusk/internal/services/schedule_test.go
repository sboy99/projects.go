package services

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sboy99/projects.go/tusk/internal/utils"
	"github.com/sboy99/projects.go/tusk/pkg/storage"
)

func TestNewScheduleService(t *testing.T) {
	service := NewScheduleService()
	if service == nil {
		t.Fatal("NewScheduleService() returned nil")
	}
	if service.cliService == nil {
		t.Error("NewScheduleService() cliService is nil")
	}
	if service.timerService == nil {
		t.Error("NewScheduleService() timerService is nil")
	}
	if service.logger == nil {
		t.Error("NewScheduleService() logger is nil")
	}
	if service.storage == nil {
		t.Error("NewScheduleService() storage is nil")
	}
}

func TestGetScriptPath(t *testing.T) {
	name := "test-schedule"
	expected := "/usr/local/bin/test-schedule.sh"
	got := getScriptPath(name)
	if got != expected {
		t.Errorf("getScriptPath(%q) = %q, want %q", name, got, expected)
	}
}

func TestGetServiceFilePath(t *testing.T) {
	name := "test-schedule"
	expected := "/etc/systemd/system/test-schedule.service"
	got := getServiceFilePath(name)
	if got != expected {
		t.Errorf("getServiceFilePath(%q) = %q, want %q", name, got, expected)
	}
}

func TestGetTimerFilePath(t *testing.T) {
	name := "test-schedule"
	expected := "/etc/systemd/system/test-schedule.timer"
	got := getTimerFilePath(name)
	if got != expected {
		t.Errorf("getTimerFilePath(%q) = %q, want %q", name, got, expected)
	}
}

func TestScheduleService_createSchedule(t *testing.T) {
	service := NewScheduleService()

	// Use a temporary storage to avoid affecting real data
	tmpDir := t.TempDir()
	service.storage = storage.NewStorage[Schedule](filepath.Join(tmpDir, "test-schedules.json"))

	t.Run("create schedule with name", func(t *testing.T) {
		err := service.createSchedule("test-schedule", "echo hello", "1m")
		if err != nil {
			t.Fatalf("createSchedule() error = %v", err)
		}
		if service.schedule == nil {
			t.Fatal("createSchedule() schedule is nil")
		}
		if service.schedule.Name != "test-schedule" {
			t.Errorf("createSchedule() Name = %q, want 'test-schedule'", service.schedule.Name)
		}
		if service.schedule.Command != "echo hello" {
			t.Errorf("createSchedule() Command = %q, want 'echo hello'", service.schedule.Command)
		}
		if service.schedule.Interval != "1m" {
			t.Errorf("createSchedule() Interval = %q, want '1m'", service.schedule.Interval)
		}
		if service.schedule.ID == "" {
			t.Error("createSchedule() ID is empty")
		}
	})

	t.Run("create schedule without name generates one", func(t *testing.T) {
		err := service.createSchedule("", "echo test", "30s")
		if err != nil {
			t.Fatalf("createSchedule() error = %v", err)
		}
		if service.schedule == nil {
			t.Fatal("createSchedule() schedule is nil")
		}
		if service.schedule.Name == "" {
			t.Error("createSchedule() Name is empty when name not provided")
		}
	})

	t.Run("create schedule with invalid command", func(t *testing.T) {
		err := service.createSchedule("test", "nonexistentcommand12345", "1m")
		if err == nil {
			t.Error("createSchedule() expected error for invalid command, got nil")
		}
	})

	t.Run("create schedule with invalid interval", func(t *testing.T) {
		err := service.createSchedule("test", "echo hello", "invalid")
		if err == nil {
			t.Error("createSchedule() expected error for invalid interval, got nil")
		}
	})
}

func TestScheduleService_getScriptContent(t *testing.T) {
	service := NewScheduleService()
	service.schedule = &Schedule{
		Command: "echo 'Hello, World!'",
	}

	content := service.getScriptContent()
	if content == "" {
		t.Error("getScriptContent() returned empty string")
	}
	if !containsHelper(content, "#!/bin/bash") {
		t.Error("getScriptContent() should start with shebang")
	}
	if !containsHelper(content, "echo 'Hello, World!'") {
		t.Error("getScriptContent() should contain the command")
	}
}

func TestScheduleService_getServiceContent(t *testing.T) {
	service := NewScheduleService()
	service.schedule = &Schedule{
		Name:           "test-schedule",
		ScriptFilePath: "/usr/local/bin/test-schedule.sh",
	}

	content := service.getServiceContent()
	if content == "" {
		t.Error("getServiceContent() returned empty string")
	}
	if !containsHelper(content, "[Unit]") {
		t.Error("getServiceContent() should contain [Unit] section")
	}
	if !containsHelper(content, "[Service]") {
		t.Error("getServiceContent() should contain [Service] section")
	}
	if !containsHelper(content, "test-schedule") {
		t.Error("getServiceContent() should contain schedule name")
	}
	if !containsHelper(content, "/usr/local/bin/test-schedule.sh") {
		t.Error("getServiceContent() should contain script file path")
	}
}

func TestScheduleService_getTimerContent(t *testing.T) {
	service := NewScheduleService()
	service.schedule = &Schedule{
		Name:     "test-schedule",
		Interval: "1m",
	}

	content := service.getTimerContent()
	if content == "" {
		t.Error("getTimerContent() returned empty string")
	}
	if !containsHelper(content, "[Unit]") {
		t.Error("getTimerContent() should contain [Unit] section")
	}
	if !containsHelper(content, "[Timer]") {
		t.Error("getTimerContent() should contain [Timer] section")
	}
	if !containsHelper(content, "test-schedule") {
		t.Error("getTimerContent() should contain schedule name")
	}
}

func TestScheduleService_getScheduleByName(t *testing.T) {
	service := NewScheduleService()

	// Use a temporary storage
	tmpDir := t.TempDir()
	service.storage = storage.NewStorage[Schedule](filepath.Join(tmpDir, "test-schedules.json"))

	// Create a test schedule
	testSchedule := Schedule{
		ID:       utils.GenerateUUID(),
		Name:     "test-schedule",
		Command:  "echo hello",
		Interval: "1m",
	}
	service.storage.Upsert(testSchedule.ID, testSchedule)

	t.Run("find existing schedule", func(t *testing.T) {
		schedule, err := service.getScheduleByName("test-schedule")
		if err != nil {
			t.Fatalf("getScheduleByName() error = %v", err)
		}
		if schedule == nil {
			t.Fatal("getScheduleByName() returned nil for existing schedule")
		}
		if schedule.Name != "test-schedule" {
			t.Errorf("getScheduleByName() Name = %q, want 'test-schedule'", schedule.Name)
		}
	})

	t.Run("find non-existent schedule", func(t *testing.T) {
		schedule, err := service.getScheduleByName("nonexistent")
		if err != nil {
			t.Fatalf("getScheduleByName() error = %v", err)
		}
		if schedule != nil {
			t.Error("getScheduleByName() expected nil for non-existent schedule")
		}
	})
}

func TestScheduleService_saveScheduleToStorage(t *testing.T) {
	service := NewScheduleService()

	// Use a temporary storage
	tmpDir := t.TempDir()
	service.storage = storage.NewStorage[Schedule](filepath.Join(tmpDir, "test-schedules.json"))

	service.schedule = &Schedule{
		ID:       utils.GenerateUUID(),
		Name:     "test-schedule",
		Command:  "echo hello",
		Interval: "1m",
	}

	err := service.saveScheduleToStorage()
	if err != nil {
		t.Fatalf("saveScheduleToStorage() error = %v", err)
	}

	// Verify it was saved
	saved, err := service.storage.Read(service.schedule.ID)
	if err != nil {
		t.Fatalf("storage.Read() error = %v", err)
	}
	if saved.Name != service.schedule.Name {
		t.Errorf("saveScheduleToStorage() saved Name = %q, want %q", saved.Name, service.schedule.Name)
	}
}

func TestScheduleService_deleteScheduleFromStorage(t *testing.T) {
	service := NewScheduleService()

	// Use a temporary storage
	tmpDir := t.TempDir()
	service.storage = storage.NewStorage[Schedule](filepath.Join(tmpDir, "test-schedules.json"))

	service.schedule = &Schedule{
		ID:       utils.GenerateUUID(),
		Name:     "test-schedule",
		Command:  "echo hello",
		Interval: "1m",
	}

	// Save first
	service.storage.Upsert(service.schedule.ID, *service.schedule)

	// Then delete
	err := service.deleteScheduleFromStorage()
	if err != nil {
		t.Fatalf("deleteScheduleFromStorage() error = %v", err)
	}

	// Verify it was deleted
	_, err = service.storage.Read(service.schedule.ID)
	if err == nil {
		t.Error("deleteScheduleFromStorage() schedule should be deleted")
	}
}

func TestScheduleService_List(t *testing.T) {
	service := NewScheduleService()

	// Use a temporary storage
	tmpDir := t.TempDir()
	service.storage = storage.NewStorage[Schedule](filepath.Join(tmpDir, "test-schedules.json"))

	t.Run("list empty schedules", func(t *testing.T) {
		err := service.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
	})

	t.Run("list with schedules", func(t *testing.T) {
		// Add a test schedule
		testSchedule := Schedule{
			ID:        utils.GenerateUUID(),
			Name:      "test-schedule",
			Command:   "echo hello",
			Interval:  "1m",
			StartTime: time.Now(),
		}
		service.storage.Upsert(testSchedule.ID, testSchedule)

		err := service.List()
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
	})
}

func TestSchedule_Struct(t *testing.T) {
	schedule := Schedule{
		ID:              "test-id",
		Name:            "test-name",
		Command:         "echo hello",
		Interval:        "1m",
		StartTime:       time.Now(),
		ScriptFilePath:  "/path/to/script.sh",
		ServiceFilePath: "/path/to/service.service",
		TimerFilePath:   "/path/to/timer.timer",
	}

	if schedule.ID != "test-id" {
		t.Errorf("Schedule.ID = %q, want 'test-id'", schedule.ID)
	}
	if schedule.Name != "test-name" {
		t.Errorf("Schedule.Name = %q, want 'test-name'", schedule.Name)
	}
	if schedule.Command != "echo hello" {
		t.Errorf("Schedule.Command = %q, want 'echo hello'", schedule.Command)
	}
}

// Helper function
func containsHelper(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Note: Testing Create() and Delete() fully would require:
// 1. Mocking CLIService to avoid actual systemd commands
// 2. Mocking file system operations
// 3. Running with sudo privileges
//
// These are better suited for integration tests.
// The unit tests above cover the core logic and helper functions.
