package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostToggle(t *testing.T) {
	t.Run("successfully toggles bulb state", func(t *testing.T) {
		env := setupTestEnv(t)

		body := bytes.NewBufferString(`{"reason":"testing toggle"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", body)
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "192.168.1.1:12345"

		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			State      bool   `json:"state"`
			Reason     string `json:"reason"`
			CreatedAt  string `json:"created_at"`
			CooldownMs int64  `json:"cooldown_ms"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.True(t, res.State)
		assert.Equal(t, "testing toggle", res.Reason)
		assert.NotEmpty(t, res.CreatedAt)
		assert.Greater(t, res.CooldownMs, int64(0))
		assert.True(t, env.engine.GetState())
	})

	t.Run("accepts empty payload", func(t *testing.T) {
		env := setupTestEnv(t)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req.RemoteAddr = "192.168.1.2:12345"

		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("rejects reason exceeding 60 chars with 400", func(t *testing.T) {
		env := setupTestEnv(t)

		longReason := strings.Repeat("a", 61)
		body := bytes.NewBufferString(`{"reason":"` + longReason + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", body)
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "192.168.1.3:12345"

		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("enforces cooldown for same IP with 429", func(t *testing.T) {
		env := setupTestEnv(t)

		body := bytes.NewBufferString(`{"reason":"first"}`)
		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", body)
		req1.RemoteAddr = "192.168.1.4:12345"
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		body2 := bytes.NewBufferString(`{"reason":"second"}`)
		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", body2)
		req2.RemoteAddr = "192.168.1.4:12345"
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
	})

	t.Run("allows different IP during active cooldown", func(t *testing.T) {
		env := setupTestEnv(t)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", bytes.NewBufferString(`{}`))
		req1.RemoteAddr = "10.0.0.1:12345"
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", bytes.NewBufferString(`{}`))
		req2.RemoteAddr = "10.0.0.2:12345"
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusOK, rec2.Code)
	})

	t.Run("enforces cooldown when second request presents issued device_id cookie", func(t *testing.T) {
		env := setupTestEnv(t)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", bytes.NewBufferString(`{}`))
		req1.RemoteAddr = "192.168.1.10:12345"
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		setCookie := rec1.Header().Get("Set-Cookie")
		require.Contains(t, setCookie, "device_id=")

		cookieParts := strings.Split(setCookie, ";")
		cookieValue := cookieParts[0]

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", bytes.NewBufferString(`{}`))
		req2.RemoteAddr = "192.168.1.10:12345"
		req2.Header.Set("Cookie", cookieValue)
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
	})

	t.Run("allows distinct devices on same IP to toggle independently", func(t *testing.T) {
		env := setupTestEnv(t)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", bytes.NewBufferString(`{}`))
		req1.RemoteAddr = "192.168.1.50:12345"
		req1.Header.Set("Cookie", "device_id=device-A")
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", bytes.NewBufferString(`{}`))
		req2.RemoteAddr = "192.168.1.50:12345"
		req2.Header.Set("Cookie", "device_id=device-B")
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusOK, rec2.Code)
	})

	t.Run("returns 500 when database insertion fails", func(t *testing.T) {
		env := setupTestEnv(t)
		_ = env.db.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", bytes.NewBufferString(`{}`))
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
