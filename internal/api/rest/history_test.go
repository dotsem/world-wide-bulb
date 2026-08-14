package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHistory(t *testing.T) {
	t.Run("returns sanitized history items without ip_hash", func(t *testing.T) {
		env := setupTestEnv(t)

		reqToggle := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		reqToggle.RemoteAddr = "192.168.1.10:12345"
		recToggle := httptest.NewRecorder()
		env.router.ServeHTTP(recToggle, reqToggle)
		require.Equal(t, http.StatusOK, recToggle.Code)

		reqHistory := httptest.NewRequest(http.MethodGet, "/api/v1/history", nil)
		recHistory := httptest.NewRecorder()
		env.router.ServeHTTP(recHistory, reqHistory)

		assert.Equal(t, http.StatusOK, recHistory.Code)

		var rawMap map[string][]map[string]any
		err := json.Unmarshal(recHistory.Body.Bytes(), &rawMap)
		require.NoError(t, err)

		toggles := rawMap["toggles"]
		require.Len(t, toggles, 1)

		item := toggles[0]
		assert.Equal(t, true, item["state"])
		assert.NotEmpty(t, item["created_at"])
		assert.Nil(t, item["ip_hash"])
	})

	t.Run("returns 500 when database query fails", func(t *testing.T) {
		env := setupTestEnv(t)
		_ = env.db.Close()

		reqHistory := httptest.NewRequest(http.MethodGet, "/api/v1/history", nil)
		recHistory := httptest.NewRecorder()
		env.router.ServeHTTP(recHistory, reqHistory)

		assert.Equal(t, http.StatusInternalServerError, recHistory.Code)
	})
}
