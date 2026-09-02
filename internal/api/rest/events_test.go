package rest_test

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamEvents(t *testing.T) {
	t.Run("sets required SSE headers and sends initial state", func(t *testing.T) {
		env := setupTestEnv(t)

		server := httptest.NewServer(env.router)
		defer server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/api/v1/events", nil)
		require.NoError(t, err)

		resp, err := server.Client().Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "text/event-stream")
		assert.Equal(t, "no-cache", resp.Header.Get("Cache-Control"))
		assert.Equal(t, "keep-alive", resp.Header.Get("Connection"))
		assert.Equal(t, "no", resp.Header.Get("X-Accel-Buffering"))

		reader := bufio.NewReader(resp.Body)

		line1, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Equal(t, "event:message\n", strings.TrimSpace(line1)+"\n")

		line2, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Contains(t, line2, `"state":false`)
	})

	t.Run("receives broadcasted events via SSE stream", func(t *testing.T) {
		env := setupTestEnv(t)

		server := httptest.NewServer(env.router)
		defer server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/api/v1/events", nil)
		require.NoError(t, err)

		resp, err := server.Client().Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		reader := bufio.NewReader(resp.Body)

		// Discard initial state lines & trailing empty line
		_, _ = reader.ReadString('\n')
		_, _ = reader.ReadString('\n')
		_, _ = reader.ReadString('\n')

		env.broker.Broadcast([]byte(`{"event":"state_changed","state":true}`))

		eventLine, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Equal(t, "event:message\n", strings.TrimSpace(eventLine)+"\n")

		dataLine, err := reader.ReadString('\n')
		require.NoError(t, err)
		assert.Contains(t, dataLine, `"state":true`)
	})
}
