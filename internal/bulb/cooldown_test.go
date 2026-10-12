package bulb

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCooldown(t *testing.T) {
	t.Run("allows toggle when IP has no prior record", func(t *testing.T) {
		c := NewCooldown(100 * time.Millisecond)
		assert.True(t, c.CanToggle("ip_fresh"))
	})

	t.Run("blocks toggle immediately after recording", func(t *testing.T) {
		c := NewCooldown(100 * time.Millisecond)
		c.Record("ip_1")

		assert.False(t, c.CanToggle("ip_1"))
	})

	t.Run("allows toggle after cooldown period expires", func(t *testing.T) {
		cooldownPeriod := 10 * time.Millisecond
		c := NewCooldown(cooldownPeriod)

		c.Record("ip_1")
		assert.False(t, c.CanToggle("ip_1"))

		time.Sleep(cooldownPeriod + 5*time.Millisecond)
		assert.True(t, c.CanToggle("ip_1"))
	})

	t.Run("isolates cooldown between different IPs after global debounce", func(t *testing.T) {
		c := NewCooldown(100 * time.Millisecond)

		c.Record("ip_1")
		assert.False(t, c.CanToggle("ip_1"))
		time.Sleep(15 * time.Millisecond)
		assert.True(t, c.CanToggle("ip_2"))
	})

	t.Run("atomic CheckAndRecord prevents concurrent double execution", func(t *testing.T) {
		c := NewCooldown(100 * time.Millisecond)
		assert.True(t, c.CheckAndRecord("ip_atomic"))
		assert.False(t, c.CheckAndRecord("ip_atomic"))
	})

	t.Run("concurrent access is thread-safe", func(_ *testing.T) {
		c := NewCooldown(50 * time.Millisecond)
		var wg sync.WaitGroup

		for range 50 {
			wg.Add(2)
			ip := "ip_shared"

			go func() {
				defer wg.Done()
				c.Record(ip)
			}()

			go func() {
				defer wg.Done()
				_ = c.CanToggle(ip)
			}()
		}

		wg.Wait()
	})

	t.Run("CheckAndRecord cleans up expired history entries when maxHistoryEntries exceeded", func(t *testing.T) {
		c := NewCooldown(1 * time.Millisecond)
		past := time.Now().Add(-5 * time.Minute)
		c.mu.Lock()
		for i := range 502 {
			c.actors[string(rune(i))] = &actorRecord{
				lastToggle:    past,
				cooldownUntil: past,
			}
		}
		c.mu.Unlock()

		assert.True(t, c.CheckAndRecord("ip_trigger_cleanup"))
	})

	t.Run("Record cleans up expired history entries when maxHistoryEntries exceeded", func(t *testing.T) {
		c := NewCooldown(1 * time.Millisecond)
		past := time.Now().Add(-5 * time.Minute)
		c.mu.Lock()
		for i := range 502 {
			c.actors[string(rune(i))] = &actorRecord{
				lastToggle:    past,
				cooldownUntil: past,
			}
		}
		c.mu.Unlock()

		c.Record("ip_trigger_cleanup")
		c.mu.Lock()
		count := len(c.actors)
		c.mu.Unlock()
		assert.Equal(t, 1, count)
	})

	t.Run("enforces global debounce across different actors", func(t *testing.T) {
		c := NewCooldown(100 * time.Millisecond)
		c.Record("actor_1")

		remaining, err := c.CheckCanToggle("actor_2")
		assert.ErrorIs(t, err, ErrGlobalDebounce)
		assert.Greater(t, remaining, time.Duration(0))

		time.Sleep(15 * time.Millisecond)

		_, err = c.CheckCanToggle("actor_2")
		assert.NoError(t, err)
	})

	t.Run("triggers 60s rapid penalty on 3rd consecutive rapid toggle", func(t *testing.T) {
		c := NewCooldown(10 * time.Millisecond)

		c.Record("clockwork_bot")
		time.Sleep(11 * time.Millisecond)

		c.Record("clockwork_bot")
		time.Sleep(11 * time.Millisecond)

		nextCooldown := c.Record("clockwork_bot")
		assert.Equal(t, 60*time.Millisecond, nextCooldown)

		remaining, err := c.CheckCanToggle("clockwork_bot")
		assert.ErrorIs(t, err, ErrRapidPenalty)
		assert.Greater(t, remaining, 10*time.Millisecond)
	})

	t.Run("resets rapid streak when actor waits past rapid window", func(t *testing.T) {
		c := NewCooldown(10 * time.Millisecond)

		c.Record("slow_actor")
		time.Sleep(20 * time.Millisecond)

		c.Record("slow_actor")
		c.mu.Lock()
		streak := c.actors["slow_actor"].rapidStreak
		c.mu.Unlock()
		assert.Equal(t, 1, streak)
	})

	t.Run("enforces velocity limit after 10 toggles within rolling window", func(t *testing.T) {
		c := NewCooldown(1 * time.Millisecond)

		for range 10 {
			c.Record("velocity_actor")
			time.Sleep(2 * time.Millisecond)
		}

		remaining, err := c.CheckCanToggle("velocity_actor")
		assert.ErrorIs(t, err, ErrVelocityLimit)
		assert.Greater(t, remaining, 3*time.Minute)
	})

	t.Run("GetRemainingCooldown returns 0 for unrecorded IP", func(t *testing.T) {
		c := NewCooldown(100 * time.Millisecond)
		assert.Equal(t, time.Duration(0), c.GetRemainingCooldown("unknown_ip"))
	})

	t.Run("GetRemainingCooldown returns remaining duration for active cooldown", func(t *testing.T) {
		c := NewCooldown(100 * time.Millisecond)
		c.Record("ip_active")
		remaining := c.GetRemainingCooldown("ip_active")
		assert.Greater(t, remaining, time.Duration(0))
		assert.LessOrEqual(t, remaining, 100*time.Millisecond)
	})

	t.Run("GetRemainingCooldown clamps to 0 when cooldown expired", func(t *testing.T) {
		c := NewCooldown(10 * time.Millisecond)
		c.Record("ip_expired")
		time.Sleep(15 * time.Millisecond)
		assert.Equal(t, time.Duration(0), c.GetRemainingCooldown("ip_expired"))
	})
}
