// Package main demonstrates how to consume the World Wide Bulb SSE stream using tmaxmax/go-sse.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmaxmax/go-sse"
)

// StateEvent represents the payload dispatched on state changes.
type StateEvent struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	State     bool   `json:"state"`
	Reason    string `json:"reason,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// ReasonEvent represents the payload dispatched on reason updates.
type ReasonEvent struct {
	Type   string `json:"type"`
	ID     int64  `json:"id"`
	Reason string `json:"reason"`
}

func main() {
	url := os.Getenv("WWB_URL")
	if url == "" {
		url = "https://wwb.dotsem.be/api/v1/events"
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create request: %v\n", err)
		os.Exit(1)
	}

	conn := sse.NewConnection(req)

	conn.SubscribeEvent("state_changed", func(event sse.Event) {
		var state StateEvent
		if err := json.Unmarshal([]byte(event.Data), &state); err == nil {
			stateStr := "OFF"
			if state.State {
				stateStr = "ON"
			}
			fmt.Printf("State: %s (ID: %d, Reason: %q)\n", stateStr, state.ID, state.Reason)
		}
	})

	conn.SubscribeEvent("reason_updated", func(event sse.Event) {
		var reason ReasonEvent
		if err := json.Unmarshal([]byte(event.Data), &reason); err == nil {
			fmt.Printf("Reason: %s (Toggle ID: %d)\n", reason.Reason, reason.ID)
		}
	})

	fmt.Printf("Connecting to %s...\n", url)
	if err := conn.Connect(); err != nil && ctx.Err() == nil {
		fmt.Fprintf(os.Stderr, "connection error: %v\n", err)
	}
}
