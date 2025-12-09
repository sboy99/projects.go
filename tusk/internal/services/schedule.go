package services

import (
	"fmt"
	"os"
	"sync"
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
	ScriptFilePath  string    `json:"scriptFilePath"`
	ServiceFilePath string    `json:"serviceFilePath"`
	TimerFilePath   string    `json:"timerFilePath"`
}

type ScheduleService struct {
	schedule     *Schedule
	cliService   *CLIService
	timerService *TimerService
	logger       *logger.Logger
	storage      *storage.Storage[Schedule]
}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{
		schedule:     &Schedule{},
		cliService:   NewCLIService(),
		timerService: NewTimerService(),
		logger:       logger.NewLogger("ScheduleService: "),
		storage:      storage.NewStorage[Schedule]("./data/schedules.json"),
	}
}

func getScriptPath(name string) string {
	return fmt.Sprintf("/usr/local/bin/%s.sh", name)
}

func getServiceFilePath(name string) string {
	return fmt.Sprintf("/etc/systemd/system/%s.service", name)
}

func getTimerFilePath(name string) string {
	return fmt.Sprintf("/etc/systemd/system/%s.timer", name)
}

func (s *ScheduleService) Create(name, command, interval string) error {
	if err := s.createSchedule(name, command, interval); err != nil {
		return err
	}

	if err := s.goCreateFiles(); err != nil {
		return err
	}
	if err := s.goEnableAndStartTimerAndService(); err != nil {
		return err
	}

	s.logger.Info("==> Reloading systemd...")
	if err := s.reloadSystemd(); err != nil {
		s.logger.Error("failed to reload systemd: %v", err)
		return err
	}

	// Use Upsert with ID as key to create or update
	if err := s.storage.Upsert(s.schedule.ID, *s.schedule); err != nil {
		s.logger.Error("failed to upsert schedule: %v", err)
		return err
	}

	s.logger.Success("scheduled task %s successfully", s.schedule.Name)
	return nil
}

func (s *ScheduleService) List() error {
	schedules, err := s.storage.ReadAll()
	if err != nil {
		s.logger.Error("failed to get all schedules: %v", err)
		return err
	}
	if len(schedules) == 0 {
		s.logger.Info("no scheduled tasks found")
		return nil
	}
	headers, rows := utils.TransformToTableData(schedules)
	s.logger.FormatTable(headers, rows)
	return nil
}

func (s *ScheduleService) Delete(name string) error {
	schedule, err := s.getScheduleByName(name)
	if err != nil {
		return err
	}
	if schedule == nil {
		s.logger.Error("schedule not found")
		return nil
	}
	if err := s.goDisableAndStopTimerAndService(); err != nil {
		return err
	}
	if err := s.goDeleteFiles(); err != nil {
		return err
	}
	if err := s.storage.Delete(schedule.ID); err != nil {
		return err
	}
	s.logger.Success("deleted schedule %s successfully", name)
	return nil
}

func (s *ScheduleService) createSchedule(name, command, interval string) error {
	if err := s.cliService.IsValidCommand(command, false); err != nil {
		s.logger.Info("==> Validating command: %s", command)
	}
	if err := s.cliService.IsValidCommand(command, false); err != nil {
		s.logger.Error("invalid command: %v", err)
		return err
	}

	s.logger.Info("==> Validating interval: %s", interval)
	if err := s.timerService.IsValidInterval(interval); err != nil {
		s.logger.Error("invalid interval: %v", err)
		return err
	}

	if name == "" {
		name = namegen.GenerateOne()
	}

	s.logger.Highlight("creating schedule %s", name)
	s.schedule = &Schedule{
		ID:              utils.GenerateUUID(),
		Name:            name,
		Command:         command,
		Interval:        interval,
		StartTime:       time.Now(),
		ScriptFilePath:  getScriptPath(name),
		ServiceFilePath: getServiceFilePath(name),
		TimerFilePath:   getTimerFilePath(name),
	}
	return nil
}

func (s *ScheduleService) getScriptContent() string {
	return fmt.Sprintf(`#!/bin/bash
%s
`, s.schedule.Command)
}

func (s *ScheduleService) createScriptFile() error {
	scriptPath := s.schedule.ScriptFilePath
	scriptContent := s.getScriptContent()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "tusk-script-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write content to temp file
	if _, err := tmpFile.WriteString(scriptContent); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Use sudo to copy temp file to destination
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo cp %s %s", tmpFile.Name(), scriptPath),
		StreamLog: false,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("failed to create script file: %s", result.Stderr)
	}
	return nil
}

func (s *ScheduleService) deleteScriptFile() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo rm -f %s", s.schedule.ScriptFilePath),
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

func (s *ScheduleService) giveScriptExecPermission() error {
	scriptPath := s.schedule.ScriptFilePath
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo chmod 755 %s", scriptPath),
		StreamLog: false,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("failed to set script permissions: %s", result.Stderr)
	}
	return nil
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
`, s.schedule.Name, s.schedule.ScriptFilePath)
}

func (s *ScheduleService) createServiceFile() error {
	serviceFilePath := s.schedule.ServiceFilePath
	serviceContent := s.getServiceContent()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "tusk-service-*.service")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write content to temp file
	if _, err := tmpFile.WriteString(serviceContent); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Use sudo to copy temp file to destination
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo cp %s %s", tmpFile.Name(), serviceFilePath),
		StreamLog: false,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("failed to create service file: %s", result.Stderr)
	}
	return nil
}

func (s *ScheduleService) deleteServiceFile() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo rm -f %s", s.schedule.ServiceFilePath),
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

func (s *ScheduleService) getTimerContent() string {
	return fmt.Sprintf(`[Unit]
Description=%s
`, s.schedule.Name)
}

func (s *ScheduleService) createTimerFile() error {
	timerFilePath := s.schedule.TimerFilePath
	timerContent := s.getTimerContent()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "tusk-timer-*.timer")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write content to temp file
	if _, err := tmpFile.WriteString(timerContent); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Use sudo to copy temp file to destination
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo cp %s %s", tmpFile.Name(), timerFilePath),
		StreamLog: false,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("failed to create timer file: %s", result.Stderr)
	}
	return nil
}

func (s *ScheduleService) deleteTimerFile() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo rm -f %s", s.schedule.TimerFilePath),
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

func (s *ScheduleService) enableTimer() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl enable --now %s.timer", s.schedule.Name),
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

func (s *ScheduleService) disableTimer() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl disable %s.timer", s.schedule.Name),
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

func (s *ScheduleService) startTimer() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl start %s.timer", s.schedule.Name),
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// Timer might not be started, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) stopTimer() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl stop %s.timer", s.schedule.Name),
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

func (s *ScheduleService) enableService() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl enable %s.service", s.schedule.Name),
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// Service might not be enabled, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) disableService() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl disable %s.service", s.schedule.Name),
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// Service might not be disabled, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) startService() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl start %s.service", s.schedule.Name),
		StreamLog: true,
	})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		// Service might not be enabled, which is okay
		return nil
	}
	return nil
}

func (s *ScheduleService) stopService() error {
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   fmt.Sprintf("sudo systemctl stop %s.service", s.schedule.Name),
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

func (s *ScheduleService) goCreateFiles() error {
	wg := sync.WaitGroup{}
	errChan := make(chan error)

	wg.Go(func() {
		s.logger.Info("==> Creating script...")
		if err := s.createScriptFile(); err != nil {
			s.logger.Error("failed to create script: %v", err)
			errChan <- err
		}
		if err := s.giveScriptExecPermission(); err != nil {
			s.logger.Error("failed to give script exec permission: %v", err)
			errChan <- err
		}
	})
	wg.Go(func() {
		s.logger.Info("==> Creating service file...")
		if err := s.createServiceFile(); err != nil {
			s.logger.Error("failed to create service file: %v", err)
			errChan <- err
		}
	})
	wg.Go(func() {
		s.logger.Info("==> Creating timer file...")
		if err := s.createTimerFile(); err != nil {
			s.logger.Error("failed to create timer file: %v", err)
			errChan <- err
		}
	})

	wg.Wait()
	close(errChan)
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (s *ScheduleService) goDeleteFiles() error {
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	wg.Go(func() {
		s.logger.Info("==> Deleting script file...")
		if err := s.deleteScriptFile(); err != nil {
			s.logger.Error("failed to delete script file: %v", err)
			errChan <- err
		}
	})
	wg.Go(func() {
		s.logger.Info("==> Deleting service file...")
		if err := s.deleteServiceFile(); err != nil {
			s.logger.Error("failed to delete service file: %v", err)
			errChan <- err
		}
	})
	wg.Go(func() {
		s.logger.Info("==> Deleting timer file...")
		if err := s.deleteTimerFile(); err != nil {
			s.logger.Error("failed to delete timer file: %v", err)
			errChan <- err
		}
	})
	wg.Wait()
	close(errChan)
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (s *ScheduleService) goEnableAndStartTimerAndService() error {
	wg := sync.WaitGroup{}
	errChan := make(chan error)

	wg.Go(func() {
		s.logger.Info("==> Enabling timer...")
		if err := s.enableTimer(); err != nil {
			s.logger.Error("failed to enable timer: %v", err)
			errChan <- err
		}
		s.logger.Info("==> Starting timer...")
		if err := s.startTimer(); err != nil {
			s.logger.Error("failed to start timer: %v", err)
			errChan <- err
		}
	})
	wg.Go(func() {
		s.logger.Info("==> Enabling service...")
		if err := s.enableService(); err != nil {
			s.logger.Error("failed to enable service: %v", err)
			errChan <- err
		}
		s.logger.Info("==> Starting service...")
		if err := s.startService(); err != nil {
			s.logger.Error("failed to start service: %v", err)
			errChan <- err
		}
	})

	wg.Wait()
	close(errChan)
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (s *ScheduleService) goDisableAndStopTimerAndService() error {
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	wg.Go(func() {
		s.logger.Info("==> Stopping timer...")
		if err := s.stopTimer(); err != nil {
			s.logger.Error("failed to stop timer: %v", err)
			errChan <- err
		}
		s.logger.Info("==> Disabling timer...")
		if err := s.disableTimer(); err != nil {
			s.logger.Error("failed to disable timer: %v", err)
			errChan <- err
		}
	})
	wg.Go(func() {
		s.logger.Info("==> Stopping service...")
		if err := s.stopService(); err != nil {
			s.logger.Error("failed to stop service: %v", err)
			errChan <- err
		}
		s.logger.Info("==> Disabling service...")
		if err := s.disableService(); err != nil {
			s.logger.Error("failed to disable service: %v", err)
			errChan <- err
		}
	})
	wg.Wait()
	close(errChan)

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (s *ScheduleService) getScheduleByName(name string) (*Schedule, error) {
	schedules, err := s.storage.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, schedule := range schedules {
		if schedule.Name == name {
			return &schedule, nil
		}
	}
	return nil, nil
}
