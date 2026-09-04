package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"f1schedule/internal/openf1"
)

type Snapshot struct {
	SavedAt  time.Time                   `json:"saved_at"`
	Meetings map[string][]openf1.Meeting `json:"meetings"`
	Sessions map[string][]openf1.Session `json:"sessions"`
}

func Path() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "f1schedule", "schedule.json")
}

func Save(path string, meetings map[int][]openf1.Meeting, sessions map[int][]openf1.Session) error {
	snap := Snapshot{
		SavedAt:  time.Now().UTC(),
		Meetings: make(map[string][]openf1.Meeting, len(meetings)),
		Sessions: make(map[string][]openf1.Session, len(sessions)),
	}
	if existing, err := Load(path); err == nil {
		snap.Meetings = existing.Meetings
		snap.Sessions = existing.Sessions
		if snap.Meetings == nil {
			snap.Meetings = map[string][]openf1.Meeting{}
		}
		if snap.Sessions == nil {
			snap.Sessions = map[string][]openf1.Session{}
		}
	}
	for year, list := range meetings {
		snap.Meetings[strconv.Itoa(year)] = list
	}
	for year, list := range sessions {
		snap.Sessions[strconv.Itoa(year)] = list
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("encode cache: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), "schedule-*.json")
	if err != nil {
		return fmt.Errorf("create cache temp: %w", err)
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write cache: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close cache: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace cache: %w", err)
	}
	return nil
}

func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("decode cache %s: %w", path, err)
	}
	return snap, nil
}

func (s Snapshot) MeetingsByYear() map[int][]openf1.Meeting {
	out := make(map[int][]openf1.Meeting, len(s.Meetings))
	for key, list := range s.Meetings {
		year, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		out[year] = list
	}
	return out
}

func (s Snapshot) SessionsByYear() map[int][]openf1.Session {
	out := make(map[int][]openf1.Session, len(s.Sessions))
	for key, list := range s.Sessions {
		year, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		out[year] = list
	}
	return out
}
