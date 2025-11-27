package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	// cronFieldCount represents the number of fields in a standard CRON expression
	cronFieldCount = 5
)

var (
	// cronFieldRegex validates a single CRON field pattern
	// Matches: *, number, number-range, */step, number/step, number-range/step
	// Examples: "*", "5", "0-59", "*/5", "0/5", "0-59/5"
	cronFieldRegex = regexp.MustCompile(`^(\*|([0-9]{1,2})(-[0-9]{1,2})?)(/([0-9]{1,2}))?$`)
)

// TimerService provides validation for timer intervals
type TimerService struct{}

// NewTimerService creates a new instance of TimerService
func NewTimerService() *TimerService {
	return &TimerService{}
}

// IsValidInterval validates that the interval is either a valid duration or CRON expression
// Duration format: "30s", "5m", "1h", etc.
// CRON format: "0 * * * *" (minute hour day month weekday)
func (t *TimerService) IsValidInterval(interval string) error {
	interval = strings.TrimSpace(interval)
	if interval == "" {
		return fmt.Errorf("interval cannot be empty")
	}

	if t.isValidDuration(interval) {
		return nil
	}

	if t.isValidCronExpression(interval) {
		return nil
	}

	return fmt.Errorf("interval must be duration (e.g., '30s', '5m') or CRON expression (e.g., '0 * * * *')")
}

// isValidDuration checks if the interval is a valid Go duration string
func (t *TimerService) isValidDuration(interval string) bool {
	_, err := time.ParseDuration(interval)
	return err == nil
}

// isValidCronExpression validates if the interval is a valid CRON expression
func (t *TimerService) isValidCronExpression(interval string) bool {
	fields := strings.Fields(interval)
	if len(fields) != cronFieldCount {
		return false
	}

	return t.areValidCronFields(fields)
}

// areValidCronFields checks if all fields match the CRON field pattern
func (t *TimerService) areValidCronFields(fields []string) bool {
	for _, field := range fields {
		if !cronFieldRegex.MatchString(field) {
			return false
		}
	}
	return true
}
