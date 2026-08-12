package bulb

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
	"world-wide-bulb/internal/store"
)

var (
	// ErrCooldown indicates a user tried to toggle before their cooldown expired.
	ErrCooldown = errors.New("rate limited: please wait before toggling again")
	// CooldownTime is the default duration a client must wait between toggles.
	CooldownTime = 10 * time.Second
)

// Store defines persistence operations needed by the bulb Engine.
type Store interface {
	GetLatestToggle(ctx context.Context) (store.Toggle, error)
	InsertToggle(ctx context.Context, arg store.InsertToggleParams) (store.Toggle, error)
}

// Engine manages the state and business logic of the bulb.
type Engine struct {
	mu       sync.Mutex
	state    atomic.Bool
	store    Store
	cooldown *Cooldown
}

// NewEngine initializes a new bulb Engine and hydrates state from the database.
func NewEngine(ctx context.Context, s Store) *Engine {
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

// RecordCooldown registers a cooldown entry for the given key hash.
func (e *Engine) RecordCooldown(ipHash string) {
	e.cooldown.Record(ipHash)
}

// Toggle flips the state of the bulb and records the toggle in the database
// It returns the new toggle record and an error if the toggle failed
// The toggle is rate limited based on the cooldown time
func (e *Engine) Toggle(ctx context.Context, reason string, ipHash string) (store.Toggle, error) {
	if !e.cooldown.CanToggle(ipHash) {
		return store.Toggle{}, ErrCooldown
	}

	e.mu.Lock()
	defer e.mu.Unlock()

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
