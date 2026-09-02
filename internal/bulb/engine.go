package bulb

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"world-wide-bulb/internal/store"

	"github.com/google/uuid"
)

var (
	// ErrCooldown indicates a user tried to toggle before their cooldown expired.
	ErrCooldown = errors.New("rate limited: please wait before toggling again")
	// ErrNotFound indicates a toggle record with the specified UUID was not found.
	ErrNotFound = errors.New("toggle not found")
	// ErrInvalidUUID indicates the provided string is not a valid UUID.
	ErrInvalidUUID = errors.New("invalid uuid format")
	// ErrReasonAlreadySet indicates a reason has already been submitted for this toggle.
	ErrReasonAlreadySet = errors.New("reason already set for this toggle")
	// ErrEmptyReason indicates the submitted reason is empty or whitespace only.
	ErrEmptyReason = errors.New("reason cannot be empty")
	// CooldownTime is the default duration a client must wait between toggles.
	CooldownTime = 10 * time.Second
)

// Store defines persistence operations needed by the bulb Engine.
type Store interface {
	GetLatestToggle(ctx context.Context) (store.Toggle, error)
	GetToggleByUUID(ctx context.Context, uuid string) (store.Toggle, error)
	InsertToggle(ctx context.Context, arg store.InsertToggleParams) (store.Toggle, error)
	UpdateToggleReason(ctx context.Context, arg store.UpdateToggleReasonParams) (sql.Result, error)
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

	toggle, err := s.GetLatestToggle(ctx)
	if err != nil {
		slog.Info("no previous state found, defaulted to false")
		e.state.Store(false)
		return e
	}

	e.state.Store(toggle.State)
	return e
}

// GetState returns the current state of the bulb.
func (e *Engine) GetState() bool {
	return e.state.Load()
}

// RecordCooldown registers a cooldown entry for the given key hash.
func (e *Engine) RecordCooldown(ipHash string) time.Duration {
	return e.cooldown.Record(ipHash)
}

// Toggle flips the state of the bulb and records the toggle in the database
// It returns the new toggle record and an error if the toggle failed
// The toggle is rate limited based on the cooldown time
func (e *Engine) Toggle(ctx context.Context, ipHash string) (store.Toggle, time.Duration, error) {
	if !e.cooldown.CanToggle(ipHash) {
		return store.Toggle{}, time.Duration(0), ErrCooldown
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.cooldown.CanToggle(ipHash) {
		return store.Toggle{}, time.Duration(0), ErrCooldown
	}

	newState := !e.state.Load()
	toggleUUID := uuid.NewString()

	toggle, err := e.store.InsertToggle(ctx, store.InsertToggleParams{
		Uuid:   toggleUUID,
		State:  newState,
		IpHash: ipHash,
	})
	if err != nil {
		slog.Error("failed to insert toggle into db", "error", err)
		return store.Toggle{}, time.Duration(0), err
	}

	remainingTime := e.cooldown.Record(ipHash)
	e.state.Store(newState)
	slog.Info("bulb toggled", "state", newState, "id", toggle.ID, "uuid", toggle.Uuid, "at", toggle.CreatedAt.Time.Format(time.RFC3339))

	return toggle, remainingTime, nil
}

// UpdateReason attaches or updates the reason for a toggle record by its UUID.
func (e *Engine) UpdateReason(ctx context.Context, id string, reason string) (store.Toggle, error) {
	trimmed := strings.TrimSpace(reason)
	if trimmed == "" {
		return store.Toggle{}, ErrEmptyReason
	}

	if _, err := uuid.Parse(id); err != nil {
		return store.Toggle{}, ErrInvalidUUID
	}

	existingToggle, err := e.store.GetToggleByUUID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.Toggle{}, ErrNotFound
		}
		return store.Toggle{}, err
	}

	if existingToggle.Reason.Valid && strings.TrimSpace(existingToggle.Reason.String) != "" {
		return store.Toggle{}, ErrReasonAlreadySet
	}

	nullReason := sql.NullString{String: trimmed, Valid: true}
	_, err = e.store.UpdateToggleReason(ctx, store.UpdateToggleReasonParams{
		Reason: nullReason,
		Uuid:   id,
	})
	if err != nil {
		return store.Toggle{}, err
	}

	existingToggle.Reason = nullReason
	return existingToggle, nil
}

// GetRemainingCooldown returns the remaining cooldown duration for the given IP hash.
func (e *Engine) GetRemainingCooldown(ipHash string) time.Duration {
	return e.cooldown.GetRemainingCooldown(ipHash)
}
