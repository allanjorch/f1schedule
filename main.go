package main

import (
	"fmt"
	"os"
	"time"

	"f1schedule/internal/cache"
	"f1schedule/internal/display"
	"f1schedule/internal/openf1"
	"f1schedule/internal/schedule"
)

func main() {
	now := time.Now()

	client := openf1.NewClient()
	weekend, err := schedule.ActiveWeekend(client, cache.Path(), now)
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
