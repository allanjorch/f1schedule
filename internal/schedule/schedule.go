package schedule

import (
	"fmt"
	"strings"
	"time"
)

func CircuitLocation(offset string) (*time.Location, error) {
	neg := false
	s := offset
	switch {
	case strings.HasPrefix(s, "-"):
		neg = true
		s = s[1:]
	case strings.HasPrefix(s, "+"):
		s = s[1:]
	}

	var hours, minutes, seconds int
	if _, err := fmt.Sscanf(s, "%d:%d:%d", &hours, &minutes, &seconds); err != nil {
		return nil, fmt.Errorf("parse gmt offset %q: %w", offset, err)
	}

	secs := hours*3600 + minutes*60 + seconds
	if neg {
		secs = -secs
	}

	label := formatUTCOffset(secs)
	return time.FixedZone(label, secs), nil
}

func formatUTCOffset(secs int) string {
	sign := "+"
	if secs < 0 {
		sign = "-"
		secs = -secs
	}
	h := secs / 3600
	m := (secs % 3600) / 60
	if m == 0 {
		return fmt.Sprintf("UTC%s%d", sign, h)
	}
	return fmt.Sprintf("UTC%s%d:%02d", sign, h, m)
}

func FormatCountdown(until time.Duration) string {
	if until < time.Minute {
		return "less than a minute"
	}

	days := int(until.Hours()) / 24
	hours := int(until.Hours()) % 24
	minutes := int(until.Minutes()) % 60

	var parts []string
	if days > 0 {
		noun := "days"
		if days == 1 {
			noun = "day"
		}
		parts = append(parts, fmt.Sprintf("%d %s", days, noun))
	}
	if hours > 0 {
		noun := "hours"
		if hours == 1 {
			noun = "hour"
		}
		parts = append(parts, fmt.Sprintf("%d %s", hours, noun))
	}
	if days == 0 && minutes > 0 {
		noun := "minutes"
		if minutes == 1 {
			noun = "minute"
		}
		parts = append(parts, fmt.Sprintf("%d %s", minutes, noun))
	}
	return strings.Join(parts, ", ")
}

func FormatDateTime(t time.Time) string {
	return t.Format("Mon 2 Jan 2006, 15:04 MST")
}