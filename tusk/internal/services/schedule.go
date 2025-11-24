package services

import (
	"fmt"
	"time"

	"github.com/sboy99/projects.go/tusk/internal/utils"
)

type ScheduleService struct {
	id        string
	command   string
	interval  string
	startTime time.Time
}

func NewScheduleService(command string, interval string) *ScheduleService {
	return &ScheduleService{
		id:        utils.GenerateUUID(),
		command:   command,
		interval:  interval,
		startTime: time.Now(),
	}
}

func (s *ScheduleService) Execute() error {
	fmt.Printf("Scheduling task '%s' at %s with ID '%s'\n", s.command, s.startTime, s.id)
	return nil
}

func (s *ScheduleService) GetID() string {
	return s.id
}
