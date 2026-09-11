package schedule

import (
	"testing"
	"time"
)

func TestZoneUses12Hour(t *testing.T) {
	tests := []struct {
		zone   string
		twelve bool
		known  bool
	}{
		{"Europe/Copenhagen", false, true},
		{"America/New_York", true, true},
		{"Australia/Sydney", true, true},
		{"Europe/London", false, true},
		{"Europe/Dublin", true, true},
		{"America/Sao_Paulo", false, true},
		{"America/Mexico_City", true, true},
		{"America/Toronto", true, true},
		{"Asia/Kolkata", true, true},
		{"Pacific/Auckland", true, true},
		{"UTC", false, true},
		{"Etc/UTC", false, true},
		{"Local", false, false},
	}
	for _, tc := range tests {
		twelve, ok := zoneUses12Hour(tc.zone)
		if ok != tc.known || twelve != tc.twelve {
			t.Fatalf("%q: got twelve=%v ok=%v, want twelve=%v ok=%v", tc.zone, twelve, ok, tc.twelve, tc.known)
		}
	}
}

func TestLocalIANAZoneUsesTZ(t *testing.T) {
	t.Setenv("TZ", "America/New_York")
	if got := localIANAZone(); got != "America/New_York" {
		t.Fatalf("TZ: got %q", got)
	}
	t.Setenv("TZ", ":Europe/Copenhagen")
	if got := localIANAZone(); got != "Europe/Copenhagen" {
		t.Fatalf("colon TZ: got %q", got)
	}
}

func TestDetect12HourClockFollowsTimezone(t *testing.T) {
	t.Setenv("TZ", "Europe/Copenhagen")
	if detect12HourClock() {
		t.Fatal("Europe/Copenhagen should use 24-hour time")
	}
	t.Setenv("TZ", "America/New_York")
	if !detect12HourClock() {
		t.Fatal("America/New_York should use 12-hour time")
	}
}

func TestFormatCompactRespectsClock(t *testing.T) {
	ts := time.Date(2026, 9, 4, 16, 0, 0, 0, time.UTC)

	twelve := true
	hour12Test = &twelve
	t.Cleanup(func() { hour12Test = nil })
	if got := FormatCompact(ts); got != "Fri 04 Sep  4:00 PM" {
		t.Fatalf("12-hour: got %q", got)
	}

	twelve = false
	if got := FormatCompact(ts); got != "Fri 04 Sep  16:00" {
		t.Fatalf("24-hour: got %q", got)
	}
}
