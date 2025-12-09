package executor

import (
	"github.com/sboy99/projects.go/tusk/internal/services"
	"github.com/sboy99/projects.go/tusk/pkg/logger"
)

type ScheduleExecutor struct {
	scheduleService *services.ScheduleService
	sudoService     *services.SudoService
	logger          *logger.Logger
}

func NewScheduleExecutor() *ScheduleExecutor {
	return &ScheduleExecutor{
		scheduleService: services.NewScheduleService(),
		sudoService:     services.NewSudoService(),
		logger:          logger.NewLogger("ScheduleExecutor: "),
	}
}

func (e *ScheduleExecutor) Create(name, command, interval string) {
	// Check and request sudo privileges before creating schedule
	if err := e.sudoService.RequestPrivileges(); err != nil {
		e.logger.Error("failed to obtain sudo privileges: %v", err)
		return
	}

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
	// Check and request sudo privileges before deleting schedule
	if err := e.sudoService.RequestPrivileges(); err != nil {
		e.logger.Error("failed to obtain sudo privileges: %v", err)
		return
	}

	if err := e.scheduleService.Delete(name); err != nil {
		e.logger.Error("failed to delete schedule: %v", err)
	}
}
