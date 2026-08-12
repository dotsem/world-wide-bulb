package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetState(t *testing.T) {
	t.Run("returns initial state false", func(t *testing.T) {
		env := setupTestEnv(t)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/state", nil)
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			State      bool  `json:"state"`
			CooldownMs int64 `json:"cooldown_ms"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.False(t, res.State)
		assert.Equal(t, int64(0), res.CooldownMs)
	})

	t.Run("returns active cooldown_ms after toggling", func(t *testing.T) {
		env := setupTestEnv(t)

		reqToggle := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		reqToggle.RemoteAddr = "192.168.1.100:12345"
		recToggle := httptest.NewRecorder()
		env.router.ServeHTTP(recToggle, reqToggle)
		require.Equal(t, http.StatusOK, recToggle.Code)

		setCookie := recToggle.Header().Get("Set-Cookie")
		require.Contains(t, setCookie, "device_id=")

		reqState := httptest.NewRequest(http.MethodGet, "/api/v1/state", nil)
		reqState.RemoteAddr = "192.168.1.100:12345"
		reqState.Header.Set("Cookie", strings.Split(setCookie, ";")[0])
		recState := httptest.NewRecorder()
		env.router.ServeHTTP(recState, reqState)

		assert.Equal(t, http.StatusOK, recState.Code)

		var res struct {
			State      bool  `json:"state"`
			CooldownMs int64 `json:"cooldown_ms"`
		}
		err := json.Unmarshal(recState.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.True(t, res.State)
		assert.Greater(t, res.CooldownMs, int64(0))
	})
}
