package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"

	"github.com/sergiught/work-pilot-cli/internal/command/base"
)

// wp work <amount-of-time-in-minutes>
// If no argument, start infinite tracker with no timeout
// wp work <time> --name "project A" // track time for project A
// wp schedule  // connect to gmail and add to schedule

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		defer cancel()
		os.Exit(0)
	}()

	log.SetLevel(log.DebugLevel)

	if err := base.NewCommand().ExecuteContext(ctx); err != nil {
		log.Fatal("failed to execute command", "err", err)
		os.Exit(1)
	}
}
