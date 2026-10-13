// Package bulb implements the core bulb engine and multi-tier rate limiting.
package bulb

import (
	"sync"
	"time"
)

const (
	maxHistoryEntries = 500
	velocityWindow    = 4 * time.Minute
	velocityLimit     = 10
	rapidStreakLimit  = 3
)

type actorRecord struct {
	cooldownUntil time.Time
	lastToggle    time.Time
	rapidStreak   int
	recentToggles []time.Time
	lastError     error
}

// Cooldown manages per-actor and global rate limiting, rapid streak penalties, and velocity quotas.
type Cooldown struct {
	mu                 sync.Mutex
	actors             map[string]*actorRecord
	lastGlobalToggle   time.Time
	cooldownTime       time.Duration
	globalDebounceTime time.Duration
	rapidStreakWindow  time.Duration
	rapidPenaltyTime   time.Duration
}

// NewCooldown initializes a new Cooldown manager with the given base duration.
func NewCooldown(cooldownTime time.Duration) *Cooldown {
	globalDebounce := 1 * time.Second
	if cooldownTime < 1*time.Second {
		globalDebounce = cooldownTime / 10
	}

	rapidWindow := 2 * time.Second
	rapidPenalty := 60 * time.Second
	if cooldownTime < 1*time.Second {
		rapidWindow = cooldownTime * 3
		rapidPenalty = cooldownTime * 6
	}

	return &Cooldown{
		actors:             make(map[string]*actorRecord),
		cooldownTime:       cooldownTime,
		globalDebounceTime: globalDebounce,
		rapidStreakWindow:  rapidWindow,
		rapidPenaltyTime:   rapidPenalty,
	}
}

func (c *Cooldown) pruneRecentToggles(history []time.Time, now time.Time) []time.Time {
	cutoff := now.Add(-velocityWindow)
	validIdx := -1
	for i, t := range history {
		if t.After(cutoff) {
			validIdx = i
			break
		}
	}
	if validIdx == -1 {
		return nil
	}
	return history[validIdx:]
}

func (c *Cooldown) cleanupExpiredLocked(now time.Time) {
	if len(c.actors) > maxHistoryEntries {
		for k, rec := range c.actors {
			if now.After(rec.cooldownUntil) && now.Sub(rec.lastToggle) >= velocityWindow {
				delete(c.actors, k)
			}
		}
	}
}

// CheckCanToggle checks if an actor can toggle and returns the remaining cooldown and error if blocked.
func (c *Cooldown) CheckCanToggle(ipHash string) (time.Duration, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	// user cooldown
	if rec, exists := c.actors[ipHash]; exists {
		if now.Before(rec.cooldownUntil) {
			err := rec.lastError
			if err == nil {
				err = ErrCooldown
			}
			return rec.cooldownUntil.Sub(now), err
		}

		rec.recentToggles = c.pruneRecentToggles(rec.recentToggles, now)
		if len(rec.recentToggles) >= velocityLimit {
			oldest := rec.recentToggles[0]
			remaining := oldest.Add(velocityWindow).Sub(now)
			if remaining > 0 {
				return remaining, ErrVelocityLimit
			}
		}
	}

	// global cooldown
	if !c.lastGlobalToggle.IsZero() {
		elapsedGlobal := now.Sub(c.lastGlobalToggle)
		if elapsedGlobal < c.globalDebounceTime {
			return c.globalDebounceTime - elapsedGlobal, ErrGlobalDebounce
		}
	}

	return 0, nil
}

// CanToggle reports whether the given actor has waited past their cooldown period.
func (c *Cooldown) CanToggle(ipHash string) bool {
	_, err := c.CheckCanToggle(ipHash)
	return err == nil
}

// CheckAndRecord atomically checks if an actor can toggle and records the timestamp.
func (c *Cooldown) CheckAndRecord(ipHash string) bool {
	if _, err := c.CheckCanToggle(ipHash); err != nil {
		return false
	}
	c.Record(ipHash)
	return true
}

// Record registers a successful toggle for the actor, updating rapid streaks and velocity limits.
// Returns the cooldown duration applied for the next attempt.
func (c *Cooldown) Record(ipHash string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	c.lastGlobalToggle = now

	rec, exists := c.actors[ipHash]
	if !exists {
		rec = &actorRecord{}
		c.actors[ipHash] = rec
	}

	// 1. Update rolling velocity history
	rec.recentToggles = c.pruneRecentToggles(rec.recentToggles, now)
	rec.recentToggles = append(rec.recentToggles, now)

	// 2. Calculate next cooldown based on rapid streaks
	var nextCooldown time.Duration
	var nextErr error

	if !rec.lastToggle.IsZero() {
		elapsed := now.Sub(rec.lastToggle)
		if elapsed >= c.cooldownTime && elapsed <= c.cooldownTime+c.rapidStreakWindow {
			rec.rapidStreak++
			if rec.rapidStreak >= rapidStreakLimit {
				nextCooldown = c.rapidPenaltyTime
				nextErr = ErrRapidPenalty
				rec.rapidStreak = 0
			} else {
				nextCooldown = c.cooldownTime
				nextErr = ErrCooldown
			}
		} else {
			rec.rapidStreak = 1
			nextCooldown = c.cooldownTime
			nextErr = ErrCooldown
		}
	} else {
		rec.rapidStreak = 1
		nextCooldown = c.cooldownTime
		nextErr = ErrCooldown
	}

	rec.lastToggle = now
	rec.cooldownUntil = now.Add(nextCooldown)
	rec.lastError = nextErr

	// 3. Extend cooldown if velocity limit reached on this toggle
	if len(rec.recentToggles) >= velocityLimit {
		oldest := rec.recentToggles[0]
		velocityCooldownUntil := oldest.Add(velocityWindow)
		if velocityCooldownUntil.After(rec.cooldownUntil) {
			rec.cooldownUntil = velocityCooldownUntil
			rec.lastError = ErrVelocityLimit
			nextCooldown = velocityCooldownUntil.Sub(now)
		}
	}

	c.cleanupExpiredLocked(now)
	return nextCooldown
}

// GetRemainingCooldown returns the remaining duration until the given actor can toggle again.
func (c *Cooldown) GetRemainingCooldown(ipHash string) time.Duration {
	remaining, _ := c.CheckCanToggle(ipHash)
	return remaining
}
