package executor

import (
	"github.com/sboy99/projects.go/tusk/internal/services"
	"github.com/sboy99/projects.go/tusk/pkg/logger"
)

type ScheduleExecutor struct {
	scheduleService *services.ScheduleService
	logger          *logger.Logger
}

func NewScheduleExecutor() *ScheduleExecutor {
	return &ScheduleExecutor{
		scheduleService: services.NewScheduleService(),
		logger:          logger.NewLogger("ScheduleExecutor: "),
	}
}

func (e *ScheduleExecutor) Create(name, command, interval string) {
	if err := e.scheduleService.Create(name, command, interval); err != nil {
		e.logger.Error("failed to create schedule: %v", err)
	}
}

func (e *ScheduleExecutor) List() {
	if err := e.scheduleService.List(); err != nil {
		e.logger.Error("failed to list schedules: %v", err)
	}
}

func (e *ScheduleExecutor) Delete(name string) {
	if err := e.scheduleService.Delete(name); err != nil {
		e.logger.Error("failed to delete schedule: %v", err)
	}
}
