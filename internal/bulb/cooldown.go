package bulb

import (
	"sync"
	"time"
)

type Cooldown struct {
	mu           sync.Mutex
	history      map[string]time.Time // ipHash -> lastToggleTime
	cooldownTime time.Duration
}

func NewCooldown(cooldownTime time.Duration) *Cooldown {
	return &Cooldown{
		history:      make(map[string]time.Time),
		cooldownTime: cooldownTime,
	}
}

// CanToggle checks if this specific IP has waited long enough
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
