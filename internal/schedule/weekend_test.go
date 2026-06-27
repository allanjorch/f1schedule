package schedule

import (
	"testing"
	"time"
)

func TestSessionStatus(t *testing.T) {
	start := time.Date(2026, 6, 27, 10, 30, 0, 0, time.UTC)
	end := time.Date(2026, 6, 27, 11, 30, 0, 0, time.UTC)

	if got := sessionStatus(start.Add(-time.Hour), start, end); got != StatusUpcoming {
		t.Fatalf("before start: got %q", got)
	}
	if got := sessionStatus(start.Add(time.Minute), start, end); got != StatusOngoing {
		t.Fatalf("during session: got %q", got)
	}
	if got := sessionStatus(end, start, end); got != StatusCompleted {
		t.Fatalf("at end: got %q", got)
	}
}

func TestIsRaceWeekend(t *testing.T) {
	if !isRaceWeekend("Australian Grand Prix") {
		t.Fatal("expected grand prix to be a race weekend")
	}
	if isRaceWeekend("Pre-Season Testing") {
		t.Fatal("expected testing to be excluded")
	}
}