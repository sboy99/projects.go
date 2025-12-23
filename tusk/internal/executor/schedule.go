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
		logger:          logger.NewLogger(""),
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

func (e *ScheduleExecutor) Logs(name string, follow bool) {
	// journalctl typically doesn't require sudo for reading logs
	// but we check anyway for consistency and in case permissions are needed
	if err := e.sudoService.RequestPrivileges(); err != nil {
		e.logger.Error("failed to obtain sudo privileges: %v", err)
		return
	}

	if err := e.scheduleService.Logs(name, follow); err != nil {
		e.logger.Error("failed to view logs: %v", err)
	}
}
