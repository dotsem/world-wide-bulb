package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Setenv("DB_PATH", ":memory:")
	t.Setenv("PORT", "-1")

	ctx := context.Background()
	err := run(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to start server")
}
