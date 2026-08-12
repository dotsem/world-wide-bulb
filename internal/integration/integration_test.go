package integration_test

import (
	"bytes"
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
	t.Cleanup(ts.Close)

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

	togglePayload := bytes.NewBufferString(`{"reason":"integration toggle test"}`)
	toggleResp, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", togglePayload)
	require.NoError(t, err)
	defer func() { _ = toggleResp.Body.Close() }()
	assert.Equal(t, http.StatusOK, toggleResp.StatusCode)

	var toggleRes struct {
		State     bool   `json:"state"`
		Reason    string `json:"reason"`
		CreatedAt string `json:"created_at"`
	}
	err = json.NewDecoder(toggleResp.Body).Decode(&toggleRes)
	require.NoError(t, err)
	assert.True(t, toggleRes.State)
	assert.Equal(t, "integration toggle test", toggleRes.Reason)

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
	assert.Equal(t, "integration toggle test", historyRes.Toggles[0].Reason)
}

func TestBackendIntegration_WebSocketBroadcast(t *testing.T) {
	env := setupIntegrationServer(t)
	client := env.server.Client()

	wsConn, _, err := websocket.DefaultDialer.Dial(env.wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = wsConn.Close() }()

	time.Sleep(50 * time.Millisecond)

	toggleBody := bytes.NewBufferString(`{"reason":"live ws broadcast test"}`)
	resp, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", toggleBody)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = wsConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	require.NoError(t, err)

	_, message, err := wsConn.ReadMessage()
	require.NoError(t, err)

	var broadcast struct {
		State     bool      `json:"state"`
		Reason    string    `json:"reason"`
		CreatedAt time.Time `json:"created_at"`
	}
	err = json.Unmarshal(message, &broadcast)
	require.NoError(t, err)
	assert.True(t, broadcast.State)
	assert.Equal(t, "live ws broadcast test", broadcast.Reason)
}

func TestBackendIntegration_IPCooldown(t *testing.T) {
	env := setupIntegrationServer(t)
	client := env.server.Client()

	body1 := bytes.NewBufferString(`{"reason":"first attempt"}`)
	resp1, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", body1)
	require.NoError(t, err)
	defer func() { _ = resp1.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	body2 := bytes.NewBufferString(`{"reason":"second attempt during cooldown"}`)
	resp2, err := client.Post(env.baseURL+"/api/v1/toggle", "application/json", body2)
	require.NoError(t, err)
	defer func() { _ = resp2.Body.Close() }()
	assert.Equal(t, http.StatusTooManyRequests, resp2.StatusCode)
}
