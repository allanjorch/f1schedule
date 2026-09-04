package schedule

import (
	"testing"
	"time"
)

func TestTFmtUses12Hour(t *testing.T) {
	tests := []struct {
		in     string
		twelve bool
		known  bool
	}{
		{"%r", true, true},
		{"%I:%M:%S %p", true, true},
		{`t_fmt="%I:%M:%S %p"`, true, true},
		{"%H:%M:%S", false, true},
		{"%R", false, true},
		{"%T", false, true},
		{"", false, false},
	}
	for _, tc := range tests {
		twelve, ok := tFmtUses12Hour(tc.in)
		if ok != tc.known || twelve != tc.twelve {
			t.Fatalf("%q: got twelve=%v ok=%v, want twelve=%v ok=%v", tc.in, twelve, ok, tc.twelve, tc.known)
		}
	}
}

func TestLangUses12Hour(t *testing.T) {
	if !langUses12Hour("en_US.UTF-8") {
		t.Fatal("expected en_US to use 12-hour time")
	}
	if langUses12Hour("en_GB.UTF-8") {
		t.Fatal("expected en_GB to use 24-hour time")
	}
	if langUses12Hour("da_DK.UTF-8") {
		t.Fatal("expected da_DK to use 24-hour time")
	}
	if langUses12Hour("C") {
		t.Fatal("expected C locale to use 24-hour time")
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
