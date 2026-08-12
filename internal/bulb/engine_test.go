package bulb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"world-wide-bulb/internal/store"

	"github.com/stretchr/testify/assert"
)

const (
	testIP     = "test_ip"
	testReason = "test_reason"
)

type mockStore struct {
	latest    store.Toggle
	getErr    error
	insertErr error
}

func (m *mockStore) GetLatestToggle(_ context.Context) (store.Toggle, error) {
	return m.latest, m.getErr
}

func (m *mockStore) InsertToggle(_ context.Context, arg store.InsertToggleParams) (store.Toggle, error) {
	if m.insertErr != nil {
		return store.Toggle{}, m.insertErr
	}
	return store.Toggle{
		State:  arg.State,
		Reason: arg.Reason,
		IpHash: arg.IpHash,
	}, nil
}

func TestNewEngine(t *testing.T) {
	t.Run("hydrated state from db", func(t *testing.T) {
		s := &mockStore{
			latest: store.Toggle{
				State: true,
			},
		}
		e := NewEngine(context.Background(), s)
		assert.True(t, e.GetState())
	})

	t.Run("defaulted to false on db error", func(t *testing.T) {
		s := &mockStore{
			getErr: sql.ErrNoRows,
		}
		e := NewEngine(context.Background(), s)
		assert.False(t, e.GetState())
	})
}

func TestToggle(t *testing.T) {
	ctx := context.Background()

	t.Run("successful toggle and state flip", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{latest: store.Toggle{State: false}})

		toggle, remainingTime, err := e.Toggle(ctx, testReason, testIP)
		assert.NoError(t, err)
		assert.True(t, toggle.State)
		assert.True(t, e.GetState())
		assert.Equal(t, sql.NullString{String: testReason, Valid: true}, toggle.Reason)
		assert.Equal(t, testIP, toggle.IpHash)
		assert.NotZero(t, remainingTime)
	})

	t.Run("empty reason results in invalid null string", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})

		toggle, _, err := e.Toggle(ctx, "", testIP)
		assert.NoError(t, err)
		assert.False(t, toggle.Reason.Valid)
	})

	t.Run("rate limited on rapid toggle from same IP", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})

		_, _, err := e.Toggle(ctx, testReason, testIP)
		assert.NoError(t, err)

		_, _, err = e.Toggle(ctx, testReason, testIP)
		assert.ErrorIs(t, err, ErrCooldown)
	})

	t.Run("different IP allowed during other IP cooldown", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})

		_, _, err := e.Toggle(ctx, testReason, "ip_1")
		assert.NoError(t, err)

		toggle, remainingTime, err := e.Toggle(ctx, testReason, "ip_2")
		assert.NoError(t, err)
		assert.False(t, toggle.State)
		assert.NotZero(t, remainingTime)
	})

	t.Run("db failure does not mutate state or record cooldown", func(t *testing.T) {
		s := &mockStore{insertErr: errors.New("db error")}
		e := NewEngine(ctx, s)

		_, _, err := e.Toggle(ctx, testReason, testIP)
		assert.Error(t, err)
		assert.False(t, e.GetState())

		s.insertErr = nil
		_, _, err = e.Toggle(ctx, testReason, testIP)
		assert.NoError(t, err)
	})

	t.Run("RecordCooldown registers cooldown for IP hash", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})
		e.RecordCooldown("manual_ip")

		_, _, err := e.Toggle(ctx, testReason, "manual_ip")
		assert.ErrorIs(t, err, ErrCooldown)
	})

	t.Run("GetRemainingCooldown delegates to cooldown struct", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})
		e.RecordCooldown("manual_ip")
		assert.Greater(t, e.GetRemainingCooldown("manual_ip"), int64(0))
	})
}
