package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sboy99/projects.go/tusk/internal/config"
	"github.com/sboy99/projects.go/tusk/internal/utils"
	"github.com/sboy99/projects.go/tusk/pkg/logger"
)

type ScheduleService struct {
	id           string
	command      string
	interval     string
	startTime    time.Time
	cliService   *CLIService
	timerService *TimerService
	logger       *logger.Logger
}

func NewScheduleService(command string, interval string) *ScheduleService {
	return &ScheduleService{
		id:           utils.GenerateUUID(),
		command:      command,
		interval:     interval,
		startTime:    time.Now(),
		cliService:   NewCLIService(),
		timerService: NewTimerService(),
		logger:       logger.NewLogger("  "),
	}
}

func (s *ScheduleService) Execute() {
	if err := s.cliService.IsValidCommand(s.command, false); err != nil {
		s.logger.Error("invalid command: %v", err)
		return
	}
	if err := s.timerService.IsValidInterval(s.interval); err != nil {
		s.logger.Error("invalid interval: %v", err)
		return
	}

	if err := s.createScript(); err != nil {
		s.logger.Error("failed to create script: %v", err)
	}
	if err := s.giveScriptExecPermission(); err != nil {
		s.logger.Error("failed to give script exec permission: %v", err)
	}

	if err := s.createServiceFile(); err != nil {
		s.logger.Error("failed to create service file: %v", err)
		return
	}
	if err := s.createTimerFile(); err != nil {
		s.logger.Error("failed to create timer file: %v", err)
		return
	}

	if err := s.reloadSystemd(); err != nil {
		s.logger.Error("failed to reload systemd: %v", err)
		return
	}
	if err := s.enableTimer(); err != nil {
		s.logger.Error("failed to enable timer: %v", err)
		return
	}

	s.logger.Success("scheduled task successfully")
}

func (s *ScheduleService) getServiceName() string {
	return strings.Join([]string{s.id, config.Get().AppName}, "-")
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
