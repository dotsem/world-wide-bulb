// Package rest provides HTTP handlers for the REST API endpoints.
package rest

import (
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/bulb"
	"world-wide-bulb/internal/store"
	"world-wide-bulb/internal/utils"
)

// Handler manages REST HTTP endpoints.
type Handler struct {
	queries *store.Queries
	engine  *bulb.Engine
	hub     *ws.Hub
	hasher  *utils.Hasher
}

// NewHandler creates a new REST handler with the given dependencies.
func NewHandler(q *store.Queries, e *bulb.Engine, h *ws.Hub, hasher *utils.Hasher) *Handler {
	return &Handler{
		queries: q,
		engine:  e,
		hub:     h,
		hasher:  hasher,
	}
}
