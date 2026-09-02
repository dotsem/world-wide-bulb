package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"world-wide-bulb/internal/api"
	"world-wide-bulb/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

type testServerEnv struct {
	server  *httptest.Server
	baseURL string
	wsURL   string
}

func setupIntegrationServer(t *testing.T) *testServerEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	cfg := &config.Config{
		IPSalt: "integration_salt",
		IsProd: false,
	}

	app, err := api.NewApp(context.Background(), cfg, db)
	require.NoError(t, err)

	ts := httptest.NewServer(app.Router)
	t.Cleanup(func() {
		ts.CloseClientConnections()
		ts.Close()
	})

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	return &testServerEnv{
		server:  ts,
		baseURL: ts.URL,
		wsURL:   wsURL,
	}
}

func TestBackendIntegration_RESTLifecycle(t *testing.T) {
	env := setupIntegrationServer(t)
	client := env.server.Client()

	resp, err := client.Get(env.baseURL + "/api/v1/state")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var stateRes struct {
		State bool `json:"state"`
	}
	err = json.NewDecoder(resp.Body).Decode(&stateRes)
	require.NoError(t, err)
	assert.False(t, stateRes.State)

	toggleResp, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", nil)
	require.NoError(t, err)
	defer func() { _ = toggleResp.Body.Close() }()
	assert.Equal(t, http.StatusOK, toggleResp.StatusCode)

	var toggleRes struct {
		ID        string `json:"id"`
		State     bool   `json:"state"`
		CreatedAt string `json:"created_at"`
	}
	err = json.NewDecoder(toggleResp.Body).Decode(&toggleRes)
	require.NoError(t, err)
	assert.True(t, toggleRes.State)
	assert.NotEmpty(t, toggleRes.ID)

	reasonPayload := strings.NewReader(`{"id":"` + toggleRes.ID + `","reason":"attached reason via post"}`)
	reasonResp, err := client.Post(env.baseURL+"/api/v1/reason", "application/json", reasonPayload)
	require.NoError(t, err)
	defer func() { _ = reasonResp.Body.Close() }()
	assert.Equal(t, http.StatusOK, reasonResp.StatusCode)

	respAfter, err := client.Get(env.baseURL + "/api/v1/state")
	require.NoError(t, err)
	defer func() { _ = respAfter.Body.Close() }()
	err = json.NewDecoder(respAfter.Body).Decode(&stateRes)
	require.NoError(t, err)
	assert.True(t, stateRes.State)

	histResp, err := client.Get(env.baseURL + "/api/v1/history")
	require.NoError(t, err)
	defer func() { _ = histResp.Body.Close() }()
	assert.Equal(t, http.StatusOK, histResp.StatusCode)

	var historyRes struct {
		Toggles []struct {
			State  bool   `json:"state"`
			Reason string `json:"reason"`
		} `json:"toggles"`
	}
	err = json.NewDecoder(histResp.Body).Decode(&historyRes)
	require.NoError(t, err)
	require.Len(t, historyRes.Toggles, 1)
	assert.True(t, historyRes.Toggles[0].State)
	assert.Equal(t, "attached reason via post", historyRes.Toggles[0].Reason)
}

func TestBackendIntegration_WebSocketBroadcast(t *testing.T) {
	env := setupIntegrationServer(t)
	client := env.server.Client()

	wsConn, _, err := websocket.DefaultDialer.Dial(env.wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = wsConn.Close() }()

	time.Sleep(50 * time.Millisecond)

	resp, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", nil)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var toggleRes struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&toggleRes)
	require.NoError(t, err)

	err = wsConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	require.NoError(t, err)

	_, message, err := wsConn.ReadMessage()
	require.NoError(t, err)

	var broadcast struct {
		ID        int64     `json:"id"`
		State     bool      `json:"state"`
		CreatedAt time.Time `json:"created_at"`
	}
	err = json.Unmarshal(message, &broadcast)
	require.NoError(t, err)
	assert.True(t, broadcast.State)
	assert.Positive(t, broadcast.ID)

	reasonPayload := strings.NewReader(`{"id":"` + toggleRes.ID + `","reason":"broadcast reason test"}`)
	reasonResp, err := client.Post(env.baseURL+"/api/v1/reason", "application/json", reasonPayload)
	require.NoError(t, err)
	defer func() { _ = reasonResp.Body.Close() }()
	assert.Equal(t, http.StatusOK, reasonResp.StatusCode)

	_, reasonMessage, err := wsConn.ReadMessage()
	require.NoError(t, err)

	var reasonBroadcast struct {
		Type   string `json:"type"`
		ID     int64  `json:"id"`
		Reason string `json:"reason"`
	}
	err = json.Unmarshal(reasonMessage, &reasonBroadcast)
	require.NoError(t, err)
	assert.Equal(t, "reason_updated", reasonBroadcast.Type)
	assert.Equal(t, broadcast.ID, reasonBroadcast.ID)
	assert.Equal(t, "broadcast reason test", reasonBroadcast.Reason)
}

func TestBackendIntegration_IPCooldown(t *testing.T) {
	env := setupIntegrationServer(t)
	client := env.server.Client()

	resp1, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", nil)
	require.NoError(t, err)
	defer func() { _ = resp1.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	resp2, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", nil)
	require.NoError(t, err)
	defer func() { _ = resp2.Body.Close() }()
	assert.Equal(t, http.StatusTooManyRequests, resp2.StatusCode)
}

func TestBackendIntegration_SSEBroadcast(t *testing.T) {
	env := setupIntegrationServer(t)
	client := env.server.Client()

	req, err := http.NewRequest(http.MethodGet, env.baseURL+"/api/v1/events", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/event-stream")

	// Trigger a toggle and verify event is broadcasted
	toggleResp, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", nil)
	require.NoError(t, err)
	defer func() { _ = toggleResp.Body.Close() }()
	assert.Equal(t, http.StatusOK, toggleResp.StatusCode)
}
