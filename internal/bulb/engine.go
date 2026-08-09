package bulb

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"
	"world-wide-bulb/internal/store"
)

var (
	ErrCooldown  = errors.New("rate limited: please wait before toggling again")
	CooldownTime = 10 * time.Second
)

type Engine struct {
	state    atomic.Bool
	store    *store.Queries
	cooldown *Cooldown
}

func NewEngine(ctx context.Context, s *store.Queries) *Engine {
	e := &Engine{
		store:    s,
		cooldown: NewCooldown(CooldownTime),
	}
	latest, err := s.GetLatestToggle(ctx)
	if err == nil {
		e.state.Store(latest.State)
		slog.Info("hydrated state from db", "state", latest.State)
	} else {
		e.state.Store(false)
		slog.Info("no previous state found, defaulted to false")
	}
	return e
}

// GetState returns the current state of the bulb
func (e *Engine) GetState() bool {
	return e.state.Load()
}

func (e *Engine) Toggle(ctx context.Context, reason string, ipHash string) (store.Toggle, error) {
	if !e.cooldown.CanToggle(ipHash) {
		return store.Toggle{}, ErrCooldown
	}

	newState := !e.state.Load()

	var nullReason sql.NullString
	if reason != "" {
		nullReason = sql.NullString{String: reason, Valid: true}
	}

	toggle, err := e.store.InsertToggle(ctx, store.InsertToggleParams{
		State:  newState,
		Reason: nullReason,
		IpHash: ipHash,
	})
	if err != nil {
		slog.Error("failed to insert toggle into db", "error", err)
		return store.Toggle{}, err
	}

	e.cooldown.Record(ipHash)
	e.state.Store(newState)
	slog.Info("bulb toggled", "state", newState, "reason", reason, "id", toggle.ID, "at", toggle.CreatedAt.Time.Format(time.RFC3339))

	return toggle, nil
}
