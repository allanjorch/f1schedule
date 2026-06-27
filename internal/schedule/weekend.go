package schedule

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"f1sched/internal/openf1"
)

type SessionStatus string

const (
	StatusUpcoming  SessionStatus = "upcoming"
	StatusOngoing   SessionStatus = "ongoing"
	StatusCompleted SessionStatus = "completed"
)

type SessionEvent struct {
	Label  string
	Status SessionStatus
	Start  time.Time
	End    time.Time
}

type Weekend struct {
	MeetingName string
	Location    string
	Country     string
	Circuit     string
	GmtOffset   string
	Sessions    []SessionEvent
	Ongoing     bool
}

func ActiveWeekend(client *openf1.Client, now time.Time) (*Weekend, error) {
	years := []int{now.Year()}
	if now.Month() >= time.November {
		years = append(years, now.Year()+1)
	}

	meetings := make(map[int]openf1.Meeting)
	var allSessions []openf1.Session

	for _, year := range years {
		yearMeetings, err := client.Meetings(year)
		if err != nil {
			return nil, err
		}
		for _, meeting := range yearMeetings {
			if meeting.IsCancelled || !isRaceWeekend(meeting.MeetingName) {
				continue
			}
			meetings[meeting.MeetingKey] = meeting
		}

		yearSessions, err := client.Sessions(year)
		if err != nil {
			return nil, err
		}
		allSessions = append(allSessions, yearSessions...)
	}

	type weekendCandidate struct {
		meetingKey int
		sessions   []parsedSession
	}

	byMeeting := make(map[int][]parsedSession)
	for _, session := range allSessions {
		if session.IsCancelled {
			continue
		}
		meeting, ok := meetings[session.MeetingKey]
		if !ok {
			continue
		}

		start, err := time.Parse(time.RFC3339, session.DateStart)
		if err != nil {
			continue
		}
		end, err := time.Parse(time.RFC3339, session.DateEnd)
		if err != nil {
			continue
		}

		byMeeting[session.MeetingKey] = append(byMeeting[session.MeetingKey], parsedSession{
			session: session,
			meeting: meeting,
			start:   start,
			end:     end,
		})
	}

	var ongoing, upcoming []weekendCandidate
	for meetingKey, sessions := range byMeeting {
		if len(sessions) == 0 {
			continue
		}

		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].start.Before(sessions[j].start)
		})

		first := sessions[0].start
		last := sessions[len(sessions)-1].end

		candidate := weekendCandidate{meetingKey: meetingKey, sessions: sessions}
		switch {
		case !now.Before(first) && now.Before(last):
			ongoing = append(ongoing, candidate)
		case now.Before(first):
			upcoming = append(upcoming, candidate)
		}
	}

	var chosen *weekendCandidate
	isOngoing := false

	if len(ongoing) > 0 {
		sort.Slice(ongoing, func(i, j int) bool {
			return ongoing[i].sessions[0].start.Before(ongoing[j].sessions[0].start)
		})
		chosen = &ongoing[0]
		isOngoing = true
	} else if len(upcoming) > 0 {
		sort.Slice(upcoming, func(i, j int) bool {
			return upcoming[i].sessions[0].start.Before(upcoming[j].sessions[0].start)
		})
		chosen = &upcoming[0]
	}

	if chosen == nil {
		return nil, fmt.Errorf("no upcoming F1 race weekends found")
	}

	first := chosen.sessions[0]
	weekend := &Weekend{
		MeetingName: first.meeting.MeetingName,
		Location:    first.meeting.Location,
		Country:     first.meeting.CountryName,
		Circuit:     first.session.CircuitShortName,
		GmtOffset:   first.session.GmtOffset,
		Ongoing:     isOngoing,
	}

	for _, s := range chosen.sessions {
		weekend.Sessions = append(weekend.Sessions, SessionEvent{
			Label:  sessionLabel(s.session.SessionName, s.session.SessionType),
			Status: sessionStatus(now, s.start, s.end),
			Start:  s.start,
			End:    s.end,
		})
	}

	return weekend, nil
}

type parsedSession struct {
	session openf1.Session
	meeting openf1.Meeting
	start   time.Time
	end     time.Time
}

func isRaceWeekend(name string) bool {
	return !strings.Contains(strings.ToLower(name), "testing")
}

func sessionStatus(now, start, end time.Time) SessionStatus {
	switch {
	case now.Before(start):
		return StatusUpcoming
	case !now.Before(end):
		return StatusCompleted
	default:
		return StatusOngoing
	}
}

func sessionLabel(name, sessionType string) string {
	if name == "" || strings.EqualFold(name, sessionType) {
		return sessionType
	}
	return fmt.Sprintf("%s (%s)", name, sessionType)
}

func (w *Weekend) NextUpcomingSession() *SessionEvent {
	for i := range w.Sessions {
		if w.Sessions[i].Status == StatusUpcoming {
			return &w.Sessions[i]
		}
	}
	return nil
}

func (w *Weekend) CurrentOngoingSession() *SessionEvent {
	for i := range w.Sessions {
		if w.Sessions[i].Status == StatusOngoing {
			return &w.Sessions[i]
		}
	}
	return nil
}

func StatusLabel(status SessionStatus) string {
	switch status {
	case StatusOngoing:
		return "ongoing"
	case StatusCompleted:
		return "completed"
	default:
		return "upcoming"
	}
}