package services

import (
	"fmt"
	"os"
	"time"

	"github.com/sboy99/projects.go/tusk/internal/utils"
	"github.com/sboy99/projects.go/tusk/pkg/logger"
	"github.com/sboy99/projects.go/tusk/pkg/namegen"
	"github.com/sboy99/projects.go/tusk/pkg/storage"
)

type Schedule struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Command         string    `json:"command"`
	Interval        string    `json:"interval"`
	StartTime       time.Time `json:"startTime"`
	ServiceFilePath string    `json:"serviceFilePath"`
	TimerFilePath   string    `json:"timerFilePath"`
}

type ScheduleService struct {
	id           string
	name         string
	command      string
	interval     string
	startTime    time.Time
	cliService   *CLIService
	timerService *TimerService
	logger       *logger.Logger
	storage      *storage.Storage[Schedule]
}

func NewScheduleService(name, command, interval string) *ScheduleService {
	if name == "" {
		name = namegen.GenerateOne()
	}
	return &ScheduleService{
		id:           utils.GenerateUUID(),
		name:         name,
		command:      command,
		interval:     interval,
		startTime:    time.Now(),
		cliService:   NewCLIService(),
		timerService: NewTimerService(),
		logger:       logger.NewLogger("  "),
		storage:      storage.NewStorage[Schedule]("./data/schedules.json"),
	}
}

func (s *ScheduleService) Create() {
	s.logger.Highlight("creating schedule %s", s.name)
	s.logger.Info("==> Validating command: %s", s.command)
	if err := s.cliService.IsValidCommand(s.command, false); err != nil {
		s.logger.Error("invalid command: %v", err)
		return
	}
	s.logger.Info("==> Validating interval: %s", s.interval)
	if err := s.timerService.IsValidInterval(s.interval); err != nil {
		s.logger.Error("invalid interval: %v", err)
		return
	}

	s.logger.Info("==> Creating script...")
	if err := s.createScript(); err != nil {
		s.logger.Error("failed to create script: %v", err)
	}
	if err := s.giveScriptExecPermission(); err != nil {
		s.logger.Error("failed to give script exec permission: %v", err)
	}

	s.logger.Info("==> Creating service file...")
	if err := s.createServiceFile(); err != nil {
		s.logger.Error("failed to create service file: %v", err)
		return
	}
	s.logger.Info("==> Creating timer file...")
	if err := s.createTimerFile(); err != nil {
		s.logger.Error("failed to create timer file: %v", err)
		return
	}
	s.logger.Info("==> Reloading systemd...")
	if err := s.reloadSystemd(); err != nil {
		s.logger.Error("failed to reload systemd: %v", err)
		return
	}
	s.logger.Info("==> Enabling timer...")
	if err := s.enableTimer(); err != nil {
		s.logger.Error("failed to enable timer: %v", err)
		return
	}

	// Use Upsert with ID as key to create or update
	if err := s.storage.Upsert(s.id, Schedule{
		ID:              s.id,
		Name:            s.name,
		Command:         s.command,
		Interval:        s.interval,
		StartTime:       s.startTime,
		ServiceFilePath: s.getServiceFilePath(),
		TimerFilePath:   s.getTimerFilePath(),
	}); err != nil {
		s.logger.Error("failed to upsert schedule: %v", err)
		return
	}

	s.logger.Success("scheduled task %s successfully", s.name)
}

func (s *ScheduleService) List() {
	schedules, err := s.storage.ReadAll()
	if err != nil {
		s.logger.Error("failed to get all schedules: %v", err)
		return
	}
	if len(schedules) == 0 {
		s.logger.Info("no scheduled tasks found")
		return
	}
	headers, rows := utils.TransformToTableData(schedules)
	s.logger.FormatTable(headers, rows)
}

func (s *ScheduleService) Delete(name string) {
	s.logger.Highlight("deleting schedule %s", name)
	
	// Find schedule by name
	schedules, err := s.storage.ReadAll()
	if err != nil {
		s.logger.Error("failed to get all schedules: %v", err)
		return
	}
	
	var scheduleToDelete *Schedule
	var scheduleID string
	for id, schedule := range schedules {
		if schedule.Name == name {
			scheduleToDelete = &schedule
			scheduleID = id
			break
		}
	}
	
	if scheduleToDelete == nil {
		s.logger.Error("schedule '%s' not found", name)
		return
	}
	
	// Stop and disable timer
	s.logger.Info("==> Stopping timer...")
	if err := s.stopTimer(scheduleToDelete.Name); err != nil {
		s.logger.Error("failed to stop timer: %v", err)
	}
	
	s.logger.Info("==> Disabling timer...")
	if err := s.disableTimer(scheduleToDelete.Name); err != nil {
		s.logger.Error("failed to disable timer: %v", err)
	}
	
	// Stop service
	s.logger.Info("==> Stopping service...")
	if err := s.stopService(scheduleToDelete.Name); err != nil {
		s.logger.Error("failed to stop service: %v", err)
	}
	
	// Delete timer file
	s.logger.Info("==> Deleting timer file...")
	if err := s.deleteTimerFile(scheduleToDelete.TimerFilePath); err != nil {
		s.logger.Error("failed to delete timer file: %v", err)
	}
	
	// Delete service file
	s.logger.Info("==> Deleting service file...")
	if err := s.deleteServiceFile(scheduleToDelete.ServiceFilePath); err != nil {
		s.logger.Error("failed to delete service file: %v", err)
	}
	
	// Delete script file
	scriptPath := fmt.Sprintf("/usr/local/bin/%s.sh", scheduleToDelete.Name)
	s.logger.Info("==> Deleting script file...")
	if err := s.deleteScriptFile(scriptPath); err != nil {
		s.logger.Error("failed to delete script file: %v", err)
	}
	
	// Reload systemd
	s.logger.Info("==> Reloading systemd...")
	if err := s.reloadSystemd(); err != nil {
		s.logger.Error("failed to reload systemd: %v", err)
	}
	
	// Delete from storage
	s.logger.Info("==> Removing from storage...")
	if err := s.storage.Delete(scheduleID); err != nil {
		s.logger.Error("failed to delete from storage: %v", err)
		return
	}
	
	s.logger.Success("schedule %s deleted successfully", name)
}

func (s *ScheduleService) getServiceName() string {
	return s.name
}

func (s *ScheduleService) getScriptPath() string {
	return fmt.Sprintf("/usr/local/bin/%s.sh", s.getServiceName())
}

func (s *ScheduleService) getScriptContent() string {
	return fmt.Sprintf(`#!/bin/bash
%s
`, s.command)
}

func (s *ScheduleService) createScript() error {
	scriptPath := s.getScriptPath()
	scriptContent := s.getScriptContent()
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return err
	}
	return nil
}

func (s *ScheduleService) giveScriptExecPermission() error {
	scriptPath := s.getScriptPath()
	if err := os.Chmod(scriptPath, 0755); err != nil {
		return err
	}
	return nil
}

func (s *ScheduleService) getServiceFilePath() string {
	return fmt.Sprintf("/etc/systemd/system/%s.service", s.getServiceName())
}

func (s *ScheduleService) getServiceContent() string {
	return fmt.Sprintf(`[Unit]
Description=%s
After=network-online.target

[Service]
Type=oneshot
ExecStart=%s
User=root
Group=root
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
`, s.getServiceName(), s.getScriptPath())
}

func (s *ScheduleService) createServiceFile() error {
	serviceFilePath := s.getServiceFilePath()
	serviceContent := s.getServiceContent()
	if err := os.WriteFile(serviceFilePath, []byte(serviceContent), 0644); err != nil {
		return err
	}
	return nil
}

func (s *ScheduleService) getTimerFilePath() string {
	return fmt.Sprintf("/etc/systemd/system/%s.timer", s.getServiceName())
}

func (s *ScheduleService) getTimerContent() string {
	return fmt.Sprintf(`[Unit]
Description=%s
`, s.getServiceName())
}

func (s *ScheduleService) createTimerFile() error {
	timerFilePath := s.getTimerFilePath()
	timerContent := s.getTimerContent()
	if err := os.WriteFile(timerFilePath, []byte(timerContent), 0644); err != nil {
		return err
	}
	return nil
}

func (s *ScheduleService) reloadSystemd() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   "sudo systemctl daemon-reload",
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		s.logger.Error("command exited with code %d: %s", result.ExitCode, result.Stderr)
		return err
	}
	return nil
}

func (s *ScheduleService) enableTimer() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl enable --now %s.timer", s.getServiceName()),
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		s.logger.Error("command exited with code %d: %s", result.ExitCode, result.Stderr)
		return err
	}
	return nil
}

func (s *ScheduleService) stopTimer(name string) error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl stop %s.timer", name),
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// Timer might not be running, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) disableTimer(name string) error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl disable %s.timer", name),
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// Timer might not be enabled, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) stopService(name string) error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl stop %s.service", name),
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// Service might not be running, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) deleteTimerFile(filePath string) error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo rm -f %s", filePath),
		StreamLog: false,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// File might not exist, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) deleteServiceFile(filePath string) error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo rm -f %s", filePath),
		StreamLog: false,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// File might not exist, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) deleteScriptFile(filePath string) error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo rm -f %s", filePath),
		StreamLog: false,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// File might not exist, which is okay
		return nil
	}
	return nil
}
