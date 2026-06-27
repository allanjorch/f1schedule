package display

import "testing"

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