package display

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"f1sched/internal/schedule"
	"f1sched/internal/term"
)

type Renderer struct {
	out  io.Writer
	term *term.Writer
}

func New(out io.Writer) *Renderer {
	return &Renderer{out: out, term: term.NewWriter()}
}

func (r *Renderer) Weekend(weekend *schedule.Weekend, now time.Time, circuitTZ *time.Location) {
	localZone, localOff := now.In(time.Local).Zone()
	circuitZone := circuitTZ.String()
	_, circuitOff := now.In(circuitTZ).Zone()
	sameZone := localOff == circuitOff

	r.printHeader(weekend, now)
	r.printTable(weekend, circuitTZ, sameZone, localZone, circuitZone)
}

func (r *Renderer) printHeader(weekend *schedule.Weekend, now time.Time) {
	const width = 52

	var badge string
	if weekend.Ongoing {
		badge = r.term.Badge(term.BgGreen+term.Bold, " ONGOING ")
	} else {
		badge = r.term.Badge(term.BgBlue+term.Bold, " UPCOMING ")
	}

	title := r.term.Bold(weekend.MeetingName)
	location := r.term.Dim(fmt.Sprintf("%s · %s", weekend.Circuit, weekend.Country))

	var nextLine string
	if next := weekend.NextUpcomingSession(); next != nil {
		countdown := schedule.FormatCountdown(next.Start.Sub(now))
		nextLine = r.term.Cyan("▸ ") + r.term.Bold("Next: ") +
			shortLabel(next.Label) + r.term.Dim(" · in ") + r.term.Yellow(countdown)
	} else if ongoing := weekend.CurrentOngoingSession(); ongoing != nil {
		nextLine = r.term.Yellow("▸ ") + r.term.Bold("Live now: ") + r.term.Yellow(shortLabel(ongoing.Label))
	}

	top := badge + " " + title
	inner := []string{top, location}
	if nextLine != "" {
		inner = append(inner, nextLine)
	}

	r.printBox(width, inner)
	fmt.Fprintln(r.out)
}

func (r *Renderer) printTable(
	weekend *schedule.Weekend,
	circuitTZ *time.Location,
	sameZone bool,
	localZone, circuitZone string,
) {
	type row struct {
		session schedule.SessionEvent
		name    string
		status  string
		local   string
		circuit string
	}

	var rows []row
	for _, session := range weekend.Sessions {
		local := formatCompact(session.Start.In(time.Local))
		circuit := formatCompact(session.Start.In(circuitTZ))
		rows = append(rows, row{
			session: session,
			name:    shortLabel(session.Label),
			status:  r.statusText(session.Status),
			local:   local,
			circuit: circuit,
		})
	}

	nameW, statusW, localW, circuitW := 7, 9, 16, 16
	for _, row := range rows {
		nameW = max(nameW, utf8.RuneCountInString(row.name))
		statusW = max(statusW, utf8.RuneCountInString(stripANSI(row.status)))
		localW = max(localW, utf8.RuneCountInString(row.local))
		if !sameZone {
			circuitW = max(circuitW, utf8.RuneCountInString(row.circuit))
		}
	}

	var header string
	var rule string
	if sameZone {
		timeHeader := fmt.Sprintf("Time (%s)", localZone)
		timeW := max(localW, utf8.RuneCountInString(timeHeader))
		header = "  " + padANSI(r.term.Dim("Session"), nameW) + "  " +
			padANSI(r.term.Dim("Status"), statusW) + "  " +
			padANSI(r.term.Dim(timeHeader), timeW)
		rule = "  " + strings.Repeat("─", nameW+statusW+timeW+4)

		fmt.Fprintln(r.out, header)
		fmt.Fprintln(r.out, r.term.Dim(rule))
		for _, row := range rows {
			line := "  " +
				padANSI(r.sessionName(row.session, row.name), nameW) + "  " +
				padANSI(row.status, statusW) + "  " +
				r.timeCell(row.session, row.local)
			fmt.Fprintln(r.out, r.highlightRow(row.session, line))
		}
		return
	}

	header = "  " + padANSI(r.term.Dim("Session"), nameW) + "  " +
		padANSI(r.term.Dim("Status"), statusW) + "  " +
		padANSI(r.term.Dim("Local ("+localZone+")"), localW) + "  " +
		padANSI(r.term.Dim("Circuit ("+circuitZone+")"), circuitW)
	rule = "  " + strings.Repeat("─", nameW+statusW+localW+circuitW+6)

	fmt.Fprintln(r.out, header)
	fmt.Fprintln(r.out, r.term.Dim(rule))
	for _, row := range rows {
		line := "  " +
			padANSI(r.sessionName(row.session, row.name), nameW) + "  " +
			padANSI(row.status, statusW) + "  " +
			padANSI(r.timeCell(row.session, row.local), localW) + "  " +
			r.timeCell(row.session, row.circuit)
		fmt.Fprintln(r.out, r.highlightRow(row.session, line))
	}
}

func (r *Renderer) printBox(width int, lines []string) {
	top := "╭" + strings.Repeat("─", width-2) + "╮"
	bottom := "╰" + strings.Repeat("─", width-2) + "╯"
	fmt.Fprintln(r.out, r.term.Dim(top))
	for _, line := range lines {
		visible := stripANSI(line)
		padding := width - 4 - utf8.RuneCountInString(visible)
		if padding < 0 {
			padding = 0
		}
		fmt.Fprintf(r.out, "%s %s%s %s\n", r.term.Dim("│"), line, strings.Repeat(" ", padding), r.term.Dim("│"))
	}
	fmt.Fprintln(r.out, r.term.Dim(bottom))
}

func (r *Renderer) statusText(status schedule.SessionStatus) string {
	switch status {
	case schedule.StatusCompleted:
		return r.term.Muted("done")
	case schedule.StatusOngoing:
		return r.term.Yellow("● live")
	default:
		return r.term.Cyan("soon")
	}
}

func (r *Renderer) sessionName(session schedule.SessionEvent, name string) string {
	switch session.Status {
	case schedule.StatusOngoing:
		return r.term.Yellow(name)
	case schedule.StatusUpcoming:
		return r.term.Cyan(name)
	default:
		return r.term.Muted(name)
	}
}

func (r *Renderer) timeCell(session schedule.SessionEvent, value string) string {
	switch session.Status {
	case schedule.StatusOngoing:
		return r.term.Yellow(value)
	case schedule.StatusUpcoming:
		return r.term.Cyan(value)
	default:
		return r.term.Muted(value)
	}
}

func (r *Renderer) highlightRow(session schedule.SessionEvent, line string) string {
	if session.Status != schedule.StatusOngoing {
		return line
	}
	if !r.term.Enabled() {
		return "▸ " + line
	}
	return r.term.Paint(term.Bold, "▸ "+line)
}

func shortLabel(label string) string {
	switch {
	case strings.HasPrefix(label, "Practice 1"):
		return "FP1"
	case strings.HasPrefix(label, "Practice 2"):
		return "FP2"
	case strings.HasPrefix(label, "Practice 3"):
		return "FP3"
	case strings.Contains(label, "Sprint Qualifying"):
		return "Sprint Quali"
	case label == "Qualifying":
		return "Qualifying"
	case label == "Sprint" || strings.HasPrefix(label, "Sprint ("):
		return "Sprint"
	case label == "Race" || strings.HasPrefix(label, "Race ("):
		return "Race"
	default:
		return label
	}
}

func formatCompact(t time.Time) string {
	return t.Format("Mon 02 Jan  15:04")
}

func stripANSI(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	skip := false
	for _, r := range s {
		if r == '\033' {
			skip = true
			continue
		}
		if skip {
			if r == 'm' {
				skip = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func padANSI(s string, width int) string {
	visible := utf8.RuneCountInString(stripANSI(s))
	if visible >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visible)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}