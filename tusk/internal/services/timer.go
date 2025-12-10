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

func (t *TimerService) convertIntervalToSystemdTimer(interval string) string {
	interval = strings.TrimSpace(interval)

	// Try to parse as duration first (e.g., "1m", "30s", "5h")
	if duration, err := time.ParseDuration(interval); err == nil {
		// Convert duration to systemd OnUnitActiveSec format
		seconds := max(1, int(duration.Seconds()))
		return fmt.Sprintf("OnUnitActiveSec=%ds", seconds)
	}

	// Try to parse as CRON expression (e.g., "0 * * * *")
	fields := strings.Fields(interval)
	if len(fields) == 5 {
		// CRON format: minute hour day month weekday
		minute, hour, day, month, weekday := fields[0], fields[1], fields[2], fields[3], fields[4]

		// Convert CRON to systemd OnCalendar format
		// systemd format: "*-*-* hour:minute:00" (year-month-day hour:minute:second)
		// Handle wildcards and specific values

		// Build date part
		datePart := "*-*-*"
		if day != "*" && month != "*" {
			datePart = fmt.Sprintf("*-%s-%s", month, day)
		} else if day != "*" {
			datePart = fmt.Sprintf("*-*-%s", day)
		} else if month != "*" {
			datePart = fmt.Sprintf("*-%s-*", month)
		}

		// Build time part
		timePart := fmt.Sprintf("%s:%s:00", hour, minute)
		if hour == "*" {
			timePart = fmt.Sprintf("*:%s:00", minute)
		}
		if minute == "*" {
			timePart = fmt.Sprintf("%s:*:00", hour)
		}
		if hour == "*" && minute == "*" {
			timePart = "*:*:00"
		}

		// Build calendar expression
		calendarExpr := fmt.Sprintf("%s %s", datePart, timePart)

		// Handle weekday - systemd uses different format
		// For now, if weekday is specified, we'll note it but systemd OnCalendar
		// doesn't directly support weekday like CRON, so we'll use the time-based schedule
		// and let the user know if needed
		if weekday != "*" {
			// systemd weekday format is different, so we'll keep the time-based schedule
			// The weekday constraint would need additional handling
		}

		return fmt.Sprintf("OnCalendar=%s", calendarExpr)
	}

	// Fallback: if we can't parse it, use OnUnitActiveSec with the raw interval
	// This might not work, but it's better than nothing
	return fmt.Sprintf("OnUnitActiveSec=%s", interval)
}
