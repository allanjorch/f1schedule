# f1schedule

A small CLI for the current Formula 1 race weekend. Run it to see which session is live or next, and the full timetable in your local time — handy when you want to know what's on without seeing results or standings.

## Try it

You need [Go](https://go.dev/) 1.24 or newer.

```bash
go run .
```

That prints the active or upcoming Grand Prix. This is the Italian Grand Prix weekend, captured on Friday 4 September 2026, with FP3 still to come:

![Italian Grand Prix weekend in the terminal: ONGOING, next session FP3, timetable through Sunday's race](screenshot.png)

Times are in your timezone. If the circuit uses a different offset, a second column shows circuit local time.

To install a binary on your `PATH`:

```bash
go build -o f1sched
```

## During a live session

The schedule comes from [OpenF1](https://openf1.org/). While a session is running, unauthenticated access is locked. After a successful run, the season timetable is saved to `~/.cache/f1sched/schedule.json`. The next time the API is locked, the app uses that cache, marks the output as cached, and warns if the file is from before this race weekend.
