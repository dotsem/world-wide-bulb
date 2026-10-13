package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostToggle(t *testing.T) {
	t.Run("successfully toggles bulb state", func(t *testing.T) {
		env := setupTestEnv(t)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "192.168.1.1:12345"

		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var res struct {
			ID         string `json:"id"`
			State      bool   `json:"state"`
			CreatedAt  string `json:"created_at"`
			CooldownMs int64  `json:"cooldown_ms"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.NotEmpty(t, res.ID)
		assert.True(t, res.State)
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

	t.Run("enforces cooldown for same IP with 429", func(t *testing.T) {
		env := setupTestEnv(t)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req1.RemoteAddr = "192.168.1.4:12345"
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req2.RemoteAddr = "192.168.1.4:12345"
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
	})

	t.Run("allows different IP during active cooldown after global debounce", func(t *testing.T) {
		env := setupTestEnvWithCooldown(t, 20*time.Millisecond)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req1.RemoteAddr = "10.0.0.1:12345"
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		time.Sleep(3 * time.Millisecond)

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req2.RemoteAddr = "10.0.0.2:12345"
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusOK, rec2.Code)
	})

	t.Run("blocks rapid toggle from different IP within global debounce window", func(t *testing.T) {
		env := setupTestEnv(t)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req1.RemoteAddr = "10.0.0.1:12345"
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req2.RemoteAddr = "10.0.0.2:12345"
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
		assert.Contains(t, rec2.Body.String(), "lamp just toggled: please wait a moment")
	})

	t.Run("enforces cooldown when second request presents issued device_id cookie", func(t *testing.T) {
		env := setupTestEnv(t)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req1.RemoteAddr = "192.168.1.10:12345"
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		setCookie := rec1.Header().Get("Set-Cookie")
		require.Contains(t, setCookie, "device_id=")

		cookieParts := strings.Split(setCookie, ";")
		cookieValue := cookieParts[0]

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req2.RemoteAddr = "192.168.1.10:12345"
		req2.Header.Set("Cookie", cookieValue)
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
	})

	t.Run("allows distinct devices on same IP to toggle independently", func(t *testing.T) {
		env := setupTestEnvWithCooldown(t, 20*time.Millisecond)

		cookieA := "device-A." + env.hasher.Sign("device-A")
		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req1.RemoteAddr = "192.168.1.50:12345"
		req1.Header.Set("Cookie", "device_id="+cookieA)
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		time.Sleep(3 * time.Millisecond)

		cookieB := "device-B." + env.hasher.Sign("device-B")
		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req2.RemoteAddr = "192.168.1.50:12345"
		req2.Header.Set("Cookie", "device_id="+cookieB)
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusOK, rec2.Code)
	})

	t.Run("rejects forged device cookies and prevents cooldown bypass", func(t *testing.T) {
		env := setupTestEnv(t)

		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req1.RemoteAddr = "192.168.1.99:12345"
		req1.Header.Set("Cookie", "device_id=forged-cookie-1")
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req2.RemoteAddr = "192.168.1.99:12345"
		req2.Header.Set("Cookie", "device_id=forged-cookie-2")
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
	})

	t.Run("escalates cooldown on 3 consecutive rapid toggles", func(t *testing.T) {
		env := setupTestEnvWithCooldown(t, 10*time.Millisecond)

		for i := range 3 {
			if i > 0 {
				time.Sleep(15 * time.Millisecond)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
			req.RemoteAddr = "192.168.1.77:12345"
			rec := httptest.NewRecorder()
			env.router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)

			if i == 2 {
				var res struct {
					CooldownMs int64 `json:"cooldown_ms"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &res)
				require.NoError(t, err)
				assert.Equal(t, int64(60), res.CooldownMs)
			}
		}

		req4 := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		req4.RemoteAddr = "192.168.1.77:12345"
		rec4 := httptest.NewRecorder()
		env.router.ServeHTTP(rec4, req4)
		assert.Equal(t, http.StatusTooManyRequests, rec4.Code)
		assert.Contains(t, rec4.Body.String(), "rapid toggle streak detected: taking a 1-minute breather")
	})

	t.Run("returns 500 when database insertion fails", func(t *testing.T) {
		env := setupTestEnv(t)
		_ = env.db.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
