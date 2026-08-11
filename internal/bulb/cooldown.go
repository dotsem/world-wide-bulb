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

	if len(c.history) > maxHistoryEntries {
		for k, t := range c.history {
			if now.Sub(t) >= c.cooldownTime {
				delete(c.history, k)
			}
		}
	}

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

// Record saves the current timestamp for this IP
func (c *Cooldown) Record(ipHash string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.history[ipHash] = time.Now()
}
