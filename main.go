package main

import (
	"fmt"
	"os"
	"time"

	"f1sched/internal/display"
	"f1sched/internal/openf1"
	"f1sched/internal/schedule"
)

func main() {
	now := time.Now()

	client := openf1.NewClient()
	weekend, err := schedule.ActiveWeekend(client, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	circuitTZ, err := schedule.CircuitLocation(weekend.GmtOffset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	display.New(os.Stdout).Weekend(weekend, now, circuitTZ)
}