package rest

import (
	"world-wide-bulb/internal/api/ws"
	"world-wide-bulb/internal/bulb"
	"world-wide-bulb/internal/store"
	"world-wide-bulb/internal/utils"
)

type Handler struct {
	queries *store.Queries
	engine  *bulb.Engine
	hub     *ws.Hub
	hasher  *utils.Hasher
}

func NewHandler(q *store.Queries, e *bulb.Engine, h *ws.Hub, hasher *utils.Hasher) *Handler {
	return &Handler{
		queries: q,
		engine:  e,
		hub:     h,
		hasher:  hasher,
	}
}

