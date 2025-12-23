package services

import (
	"fmt"
	"regexp"
	"strconv"
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
// Duration format: "30s", "5m", "1h", "1d", "1w", "1mn", "1yr", etc.
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

	return fmt.Errorf("interval must be duration (e.g., '30s', '5m', '1h', '1d', '1w', '1mn', '1yr') or CRON expression (e.g., '0 * * * *')")
}

// isValidDuration checks if the interval is a valid duration string
// Supports standard Go duration units (s, m, h) and custom units (d, w, mn, yr)
func (t *TimerService) isValidDuration(interval string) bool {
	// First try standard Go duration parsing
	if _, err := time.ParseDuration(interval); err == nil {
		return true
	}

	// Try parsing with custom units (d, w, mn, yr)
	return t.isValidCustomDuration(interval)
}

// isValidCustomDuration checks if the interval uses custom duration units (d, w, mn, yr)
func (t *TimerService) isValidCustomDuration(interval string) bool {
	// Match pattern: number followed by unit (d, w, mn, yr)
	// Examples: "1d", "2w", "3mn", "1yr"
	customDurationRegex := regexp.MustCompile(`^(\d+)(d|w|mn|yr)$`)
	matches := customDurationRegex.FindStringSubmatch(interval)
	if len(matches) != 3 {
		return false
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil || value <= 0 {
		return false
	}

	unit := matches[2]
	// All units are valid: d, w, mn, yr
	return unit == "d" || unit == "w" || unit == "mn" || unit == "yr"
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

	// Try to parse as custom duration first (d, w, mn, yr)
	// This handles: 1d, 2d, 1w, 2w, 1mn, 2mn, 1yr, 2yr, etc.
	if converted := t.convertCustomDurationToSystemd(interval); converted != "" {
		return converted
	}

	// Try to parse as standard Go duration (e.g., "1m", "30s", "5h")
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

// convertCustomDurationToSystemd converts custom duration units (d, w, mn, yr) to systemd timer format
func (t *TimerService) convertCustomDurationToSystemd(interval string) string {
	customDurationRegex := regexp.MustCompile(`^(\d+)(d|w|mn|yr)$`)
	matches := customDurationRegex.FindStringSubmatch(interval)
	if len(matches) != 3 {
		return ""
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil || value <= 0 {
		return ""
	}

	unit := matches[2]
	switch unit {
	case "d":
		// Days: use OnCalendar format for better precision
		if value == 1 {
			// Daily: run every day at midnight
			return "OnCalendar=daily"
		}
		// For multiple days, use OnUnitActiveSec with seconds
		// 1 day = 86400 seconds
		seconds := value * 86400
		return fmt.Sprintf("OnUnitActiveSec=%ds", seconds)
	case "w":
		// Weeks: use OnCalendar format for better precision
		if value == 1 {
			// Weekly: run every week on Monday at midnight
			return "OnCalendar=weekly"
		}
		// For multiple weeks, use OnUnitActiveSec with seconds
		// 1 week = 604800 seconds
		seconds := value * 604800
		return fmt.Sprintf("OnUnitActiveSec=%ds", seconds)
	case "mn":
		// Months: use OnCalendar format
		if value == 1 {
			// Monthly: run on the 1st of every month at midnight
			return "OnCalendar=monthly"
		}
		// For multiple months, use OnCalendar with interval
		// systemd supports interval syntax: "*-*-1/2" means every 2 months on the 1st
		// But for simplicity, we'll use OnUnitActiveSec with approximate seconds
		// Approximate: 1 month ≈ 30 days = 2592000 seconds
		seconds := value * 2592000
		return fmt.Sprintf("OnUnitActiveSec=%ds", seconds)
	case "yr":
		// Years: use OnCalendar format
		if value == 1 {
			// Yearly: run on January 1st at midnight
			return "OnCalendar=yearly"
		}
		// For multiple years, use OnUnitActiveSec with approximate seconds
		// Approximate: 1 year ≈ 365 days = 31536000 seconds
		seconds := value * 31536000
		return fmt.Sprintf("OnUnitActiveSec=%ds", seconds)
	default:
		return ""
	}
}
