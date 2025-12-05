package executor

import "github.com/sboy99/projects.go/tusk/internal/services"

type ScheduleExecutor struct {
	scheduleService *services.ScheduleService
}

func NewScheduleExecutor(name, command, interval string) *ScheduleExecutor {
	return &ScheduleExecutor{scheduleService: services.NewScheduleService(name, command, interval)}
}

func (e *ScheduleExecutor) Create() {
	e.scheduleService.Create()
}

func (e *ScheduleExecutor) List() {
	e.scheduleService.List()
}

func (e *ScheduleExecutor) Delete(name string) {
	e.scheduleService.Delete(name)
}
