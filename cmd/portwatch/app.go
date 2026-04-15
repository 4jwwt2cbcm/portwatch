package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// App wires together the scanner, state manager, and alert dispatcher
// into a runnable daemon loop.
type App struct {
	cfg        *config.Config
	manager    *state.Manager
	dispatcher *alert.Dispatcher
}

// NewApp constructs an App from the provided configuration.
func NewApp(cfg *config.Config) (*App, error) {
	sc := scanner.NewScanner(cfg.Ports, cfg.Host)

	store, err := state.NewStore(cfg.StateFile)
	if err != nil {
		return nil, fmt.Errorf("state store: %w", err)
	}

	notifier := alert.NewNotifier(os.Stdout)
	dispatcher := alert.NewDispatcher(notifier)
	manager := state.NewManager(sc, store, dispatcher)

	return &App{
		cfg:        cfg,
		manager:    manager,
		dispatcher: dispatcher,
	}, nil
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (a *App) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(a.cfg.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	fmt.Printf("portwatch started (interval: %ds)\n", a.cfg.IntervalSeconds)

	// Run once immediately before waiting for the first tick.
	if err := a.manager.Cycle(); err != nil {
		fmt.Fprintf(os.Stderr, "cycle error: %v\n", err)
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("portwatch shutting down")
			return
		case <-ticker.C:
			if err := a.manager.Cycle(); err != nil {
				fmt.Fprintf(os.Stderr, "cycle error: %v\n", err)
			}
		}
	}
}
