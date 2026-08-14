package bulb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"world-wide-bulb/internal/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testIP = "test_ip"

type mockResult struct {
	rows int64
	err  error
}

func (m mockResult) LastInsertId() (int64, error) { return 0, nil }
func (m mockResult) RowsAffected() (int64, error) { return m.rows, m.err }

type mockStore struct {
	latest       store.Toggle
	getErr       error
	insertErr    error
	updateErr    error
	rowsAffected int64
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
		IpHash: arg.IpHash,
		Uuid:   arg.Uuid,
	}, nil
}

func (m *mockStore) UpdateToggleReason(_ context.Context, _ store.UpdateToggleReasonParams) (sql.Result, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return mockResult{rows: m.rowsAffected}, nil
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

		toggle, remainingTime, err := e.Toggle(ctx, testIP)
		assert.NoError(t, err)
		assert.True(t, toggle.State)
		assert.True(t, e.GetState())
		assert.Equal(t, testIP, toggle.IpHash)
		assert.NotZero(t, remainingTime)
	})

	t.Run("rate limited on rapid toggle from same IP", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})

		_, _, err := e.Toggle(ctx, testIP)
		assert.NoError(t, err)

		_, _, err = e.Toggle(ctx, testIP)
		assert.ErrorIs(t, err, ErrCooldown)
	})

	t.Run("different IP allowed during other IP cooldown", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})

		_, _, err := e.Toggle(ctx, "ip_1")
		assert.NoError(t, err)

		toggle, remainingTime, err := e.Toggle(ctx, "ip_2")
		assert.NoError(t, err)
		assert.False(t, toggle.State)
		assert.NotZero(t, remainingTime)
	})

	t.Run("db failure does not mutate state or record cooldown", func(t *testing.T) {
		s := &mockStore{insertErr: errors.New("db error")}
		e := NewEngine(ctx, s)

		_, _, err := e.Toggle(ctx, testIP)
		assert.Error(t, err)
		assert.False(t, e.GetState())

		s.insertErr = nil
		_, _, err = e.Toggle(ctx, testIP)
		assert.NoError(t, err)
	})

	t.Run("RecordCooldown registers cooldown for IP hash", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})
		e.RecordCooldown("manual_ip")

		_, _, err := e.Toggle(ctx, "manual_ip")
		assert.ErrorIs(t, err, ErrCooldown)
	})

	t.Run("GetRemainingCooldown delegates to cooldown struct", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})
		e.RecordCooldown("manual_ip")
		assert.Greater(t, e.GetRemainingCooldown("manual_ip"), int64(0))
	})

	t.Run("UpdateReason attaches reason for valid UUID", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{rowsAffected: 1})
		toggle, _, err := e.Toggle(ctx, testIP)
		require.NoError(t, err)

		err = e.UpdateReason(ctx, toggle.Uuid, "new reason")
		assert.NoError(t, err)
	})

	t.Run("UpdateReason returns ErrInvalidUUID for invalid string", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{})
		err := e.UpdateReason(ctx, "invalid-uuid", "reason")
		assert.ErrorIs(t, err, ErrInvalidUUID)
	})

	t.Run("UpdateReason returns ErrNotFound when 0 rows affected", func(t *testing.T) {
		e := NewEngine(ctx, &mockStore{rowsAffected: 0})
		err := e.UpdateReason(ctx, "123e4567-e89b-12d3-a456-426614174000", "reason")
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
