package executor

import "github.com/sboy99/projects.go/tusk/internal/services"

type ScheduleExecutor struct {
	scheduleService *services.ScheduleService
}

func NewScheduleExecutor(command string, interval string) *ScheduleExecutor {
	return &ScheduleExecutor{scheduleService: services.NewScheduleService(command, interval)}
}

func (e *ScheduleExecutor) Execute() error {
	return e.scheduleService.Execute()
}

func (e *ScheduleExecutor) GetID() string {
	return e.scheduleService.GetID()
}
