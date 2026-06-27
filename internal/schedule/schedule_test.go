package schedule

import (
	"testing"
	"time"
)

func TestCircuitLocation(t *testing.T) {
	tests := []struct {
		offset string
		want   string
	}{
		{"02:00:00", "UTC+2"},
		{"-04:00:00", "UTC-4"},
		{"11:00:00", "UTC+11"},
		{"05:30:00", "UTC+5:30"},
	}

	for _, tc := range tests {
		loc, err := CircuitLocation(tc.offset)
		if err != nil {
			t.Fatalf("offset %q: %v", tc.offset, err)
		}
		if loc.String() != tc.want {
			t.Fatalf("offset %q: got %q, want %q", tc.offset, loc.String(), tc.want)
		}
	}
}

func TestFormatCountdown(t *testing.T) {
	if got := FormatCountdown(45 * time.Minute); got != "45 minutes" {
		t.Fatalf("got %q", got)
	}
	if got := FormatCountdown(26*time.Hour + 30*time.Minute); got != "1 day, 2 hours" {
		t.Fatalf("got %q", got)
	}
}