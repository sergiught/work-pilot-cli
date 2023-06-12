package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"

	"github.com/sergiught/work-pilot-cli/internal/command/base"
)

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
	}
}
