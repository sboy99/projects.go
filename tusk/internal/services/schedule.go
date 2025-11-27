package services

import (
	"fmt"
	"time"

	"github.com/sboy99/projects.go/tusk/internal/utils"
)

type ScheduleService struct {
	id         string
	command    string
	interval   string
	startTime  time.Time
	cliService *CLIService
}

func NewScheduleService(command string, interval string) *ScheduleService {
	return &ScheduleService{
		id:         utils.GenerateUUID(),
		command:    command,
		interval:   interval,
		startTime:  time.Now(),
		cliService: NewCLIService(),
	}
}

func (s *ScheduleService) Execute() error {
	fmt.Printf("Scheduling task '%s' at %s with ID '%s'\n", s.command, s.startTime.Format(time.RFC3339), s.id)

	if err := s.cliService.IsValidCommand(s.command, false); err != nil {
		return fmt.Errorf("invalid command: %w", err)
	}

	// Execute the command using CLI service with streaming enabled
	result, err := s.cliService.Execute(ExecuteOptions{
		Command:   s.command,
		StreamLog: true,
	})

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("command exited with code %d: %s", result.ExitCode, result.Stderr)
	}

	return nil
}

func (s *ScheduleService) GetID() string {
	return s.id
}
