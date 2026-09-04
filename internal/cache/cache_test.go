package cache

import (
	"path/filepath"
	"testing"

	"f1sched/internal/openf1"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "schedule.json")
	meetings := map[int][]openf1.Meeting{
		2026: {{MeetingKey: 1, MeetingName: "Italian Grand Prix"}},
	}
	sessions := map[int][]openf1.Session{
		2026: {{SessionKey: 9, SessionName: "Practice 1", MeetingKey: 1}},
	}

	if err := Save(path, meetings, sessions); err != nil {
		t.Fatal(err)
	}

	snap, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if snap.SavedAt.IsZero() {
		t.Fatal("expected saved_at")
	}
	gotMeetings := snap.MeetingsByYear()[2026]
	if len(gotMeetings) != 1 || gotMeetings[0].MeetingName != "Italian Grand Prix" {
		t.Fatalf("meetings: %+v", gotMeetings)
	}
	gotSessions := snap.SessionsByYear()[2026]
	if len(gotSessions) != 1 || gotSessions[0].SessionName != "Practice 1" {
		t.Fatalf("sessions: %+v", gotSessions)
	}
}

func TestSaveMergesYears(t *testing.T) {
	path := filepath.Join(t.TempDir(), "schedule.json")
	if err := Save(path, map[int][]openf1.Meeting{
		2026: {{MeetingKey: 1, MeetingName: "Italian Grand Prix"}},
	}, map[int][]openf1.Session{
		2026: {{SessionKey: 1, MeetingKey: 1}},
	}); err != nil {
		t.Fatal(err)
	}

	if err := Save(path, map[int][]openf1.Meeting{
		2027: {{MeetingKey: 2, MeetingName: "Australian Grand Prix"}},
	}, map[int][]openf1.Session{
		2027: {{SessionKey: 2, MeetingKey: 2}},
	}); err != nil {
		t.Fatal(err)
	}

	snap, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := snap.MeetingsByYear()[2026]; !ok {
		t.Fatal("expected 2026 meetings to be kept")
	}
	if _, ok := snap.MeetingsByYear()[2027]; !ok {
		t.Fatal("expected 2027 meetings to be added")
	}
}
