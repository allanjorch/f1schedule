package display

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"f1schedule/internal/schedule"
)

func TestShortLabel(t *testing.T) {
	tests := map[string]string{
		"Practice 1 (Practice)": "FP1",
		"Qualifying":            "Qualifying",
		"Sprint Qualifying":     "Sprint Quali",
		"Race":                  "Race",
	}
	for input, want := range tests {
		if got := shortLabel(input); got != want {
			t.Fatalf("%q: got %q, want %q", input, got, want)
		}
	}
}

func TestStripANSI(t *testing.T) {
	if got := stripANSI("\033[32mdone\033[0m"); got != "done" {
		t.Fatalf("got %q", got)
	}
}

func TestPadANSI(t *testing.T) {
	got := padANSI("\033[36msoon\033[0m", 8)
	if stripANSI(got) != "soon    " {
		t.Fatalf("got %q", stripANSI(got))
	}
}

func TestWeekendShowsCachedAndLiveSession(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	now := time.Date(2026, 9, 4, 14, 0, 0, 0, time.UTC)
	weekend := &schedule.Weekend{
		MeetingName:  "Italian Grand Prix",
		Circuit:      "Monza",
		Country:      "Italy",
		GmtOffset:    "02:00:00",
		Ongoing:      true,
		Cached:       true,
		CacheSavedAt: time.Date(2026, 9, 4, 12, 4, 0, 0, time.UTC),
		Sessions: []schedule.SessionEvent{
			{
				Label:  "Practice 1 (Practice)",
				Status: schedule.StatusOngoing,
				Start:  now.Add(-time.Hour),
				End:    now.Add(time.Hour),
			},
			{
				Label:  "Practice 2 (Practice)",
				Status: schedule.StatusUpcoming,
				Start:  now.Add(2 * time.Hour),
				End:    now.Add(3 * time.Hour),
			},
		},
	}

	var buf bytes.Buffer
	New(&buf).Weekend(weekend, now, time.UTC)
	out := buf.String()
	if !strings.Contains(out, "cached") {
		t.Fatalf("expected cached marker in:\n%s", out)
	}
	if !strings.Contains(out, "Live now: FP1") {
		t.Fatalf("expected live session in:\n%s", out)
	}
	if !strings.Contains(out, "Next: FP2") {
		t.Fatalf("expected next session in:\n%s", out)
	}
	if strings.Contains(out, "STALE CACHE") {
		t.Fatalf("did not expect a stale warning for a same-weekend cache:\n%s", out)
	}
}

func TestWeekendShowsStaleCacheWarning(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	now := time.Date(2026, 9, 4, 14, 0, 0, 0, time.UTC)
	weekend := &schedule.Weekend{
		MeetingName:  "Italian Grand Prix",
		Circuit:      "Monza",
		Country:      "Italy",
		Ongoing:      true,
		Cached:       true,
		CacheSavedAt: now.Add(-10 * 24 * time.Hour),
		Sessions: []schedule.SessionEvent{
			{
				Label:  "Practice 1 (Practice)",
				Status: schedule.StatusOngoing,
				Start:  now.Add(-time.Hour),
				End:    now.Add(time.Hour),
			},
		},
	}

	var buf bytes.Buffer
	New(&buf).Weekend(weekend, now, time.UTC)
	out := buf.String()
	if !strings.Contains(out, "STALE CACHE") {
		t.Fatalf("expected STALE CACHE badge in:\n%s", out)
	}
	if !strings.Contains(out, "cache is not from this race weekend") {
		t.Fatalf("expected stale explanation in:\n%s", out)
	}
}

func TestWeekendShowsPastWeekendWarning(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	now := time.Date(2026, 9, 4, 14, 0, 0, 0, time.UTC)
	weekend := &schedule.Weekend{
		MeetingName:  "Belgian Grand Prix",
		Circuit:      "Spa",
		Country:      "Belgium",
		Cached:       true,
		CacheSavedAt: now.Add(-10 * 24 * time.Hour),
		Sessions: []schedule.SessionEvent{
			{
				Label:  "Race",
				Status: schedule.StatusCompleted,
				Start:  now.Add(-10 * 24 * time.Hour),
				End:    now.Add(-10*24*time.Hour + 2*time.Hour),
			},
		},
	}

	var buf bytes.Buffer
	New(&buf).Weekend(weekend, now, time.UTC)
	out := buf.String()
	if !strings.Contains(out, "PAST") {
		t.Fatalf("expected PAST badge in:\n%s", out)
	}
	if !strings.Contains(out, "this race weekend has already finished") {
		t.Fatalf("expected finished warning in:\n%s", out)
	}
	if !strings.Contains(out, "cached from a previous GP") {
		t.Fatalf("expected previous-GP cache warning in:\n%s", out)
	}
}
