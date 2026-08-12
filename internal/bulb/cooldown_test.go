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

	t.Run("isolates cooldown between different IPs", func(t *testing.T) {
		c := NewCooldown(100 * time.Millisecond)

		c.Record("ip_1")
		assert.False(t, c.CanToggle("ip_1"))
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
		for i := range 502 {
			c.Record(string(rune(i)))
		}
		time.Sleep(2 * time.Millisecond)

		assert.True(t, c.CheckAndRecord("ip_trigger_cleanup"))
	})
}
