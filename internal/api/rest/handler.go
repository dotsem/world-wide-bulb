// Package rest provides HTTP handlers for the REST API endpoints.
package rest

import (
	"world-wide-bulb/internal/api/sse"
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/bulb"
	"world-wide-bulb/internal/store"
	"world-wide-bulb/internal/utils"
)

const (
	stateKey          = "state"
	eventKey          = "event"
	eventStateChanged = "state_changed"
)

// Handler manages REST HTTP endpoints.
type Handler struct {
	queries *store.Queries
	engine  *bulb.Engine
	hub     *ws.Hub
	broker  *sse.Broker
	hasher  *utils.Hasher
	isProd  bool
}

// NewHandler creates a new REST handler with the given dependencies.
func NewHandler(q *store.Queries, e *bulb.Engine, h *ws.Hub, b *sse.Broker, hasher *utils.Hasher, isProd bool) *Handler {
	return &Handler{
		queries: q,
		engine:  e,
		hub:     h,
		broker:  b,
		hasher:  hasher,
		isProd:  isProd,
	}
}
