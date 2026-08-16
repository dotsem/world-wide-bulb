// Package bulb implements the core bulb engine and per-IP cooldown rate limiting.
package bulb

import (
	"sync"
	"time"
)

const maxHistoryEntries = 500

// Cooldown manages per-IP rate limiting and history tracking.
type Cooldown struct {
	mu           sync.Mutex
	history      map[string]time.Time
	cooldownTime time.Duration
}

// NewCooldown initializes a new Cooldown manager with the given duration.
func NewCooldown(cooldownTime time.Duration) *Cooldown {
	return &Cooldown{
		history:      make(map[string]time.Time),
		cooldownTime: cooldownTime,
	}
}

func (c *Cooldown) cleanupExpiredLocked(now time.Time) {
	if len(c.history) > maxHistoryEntries {
		for k, t := range c.history {
			if now.Sub(t) >= c.cooldownTime {
				delete(c.history, k)
			}
		}
	}
}

// CheckAndRecord atomically checks if an IP can toggle and records the timestamp.
func (c *Cooldown) CheckAndRecord(ipHash string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if lastTime, exists := c.history[ipHash]; exists {
		if now.Sub(lastTime) < c.cooldownTime {
			return false
		}
	}

	c.history[ipHash] = now
	c.cleanupExpiredLocked(now)

	return true
}

// CanToggle reports whether the given IP hash has waited past the cooldown period.
func (c *Cooldown) CanToggle(ipHash string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	lastTime, exists := c.history[ipHash]
	if !exists {
		return true
	}

	return time.Since(lastTime) >= c.cooldownTime
}

// Record saves the current timestamp for this IP and prunes expired records if capacity exceeded.
// Returns the remaining time until the next toggle.
func (c *Cooldown) Record(ipHash string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	c.history[ipHash] = now
	c.cleanupExpiredLocked(now)

	return c.cooldownTime
}

// GetRemainingCooldown returns the remaining duration until the given IP hash can toggle again.
func (c *Cooldown) GetRemainingCooldown(ipHash string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	lastTime, exists := c.history[ipHash]
	if !exists {
		return time.Duration(0)
	}

	remaining := c.cooldownTime - time.Since(lastTime)
	if remaining < 0 {
		return 0
	}
	return remaining
}
