package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"world-wide-bulb/internal/store"

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

		var rawMap map[string]any
		err := json.Unmarshal(recHistory.Body.Bytes(), &rawMap)
		require.NoError(t, err)

		toggles, ok := rawMap["toggles"].([]any)
		require.True(t, ok)
		require.Len(t, toggles, 1)

		item := toggles[0].(map[string]any)
		assert.Equal(t, true, item["state"])
		assert.NotEmpty(t, item["created_at"])
		assert.Nil(t, item["ip_hash"])

		assert.Equal(t, false, rawMap["has_more"])
	})

	t.Run("supports limit and cursor pagination", func(t *testing.T) {
		env := setupTestEnv(t)
		ctx := context.Background()

		for i := range 3 {
			_, err := env.queries.InsertToggle(ctx, store.InsertToggleParams{
				Uuid:   fmt.Sprintf("uuid-%d", i),
				State:  i%2 == 0,
				IpHash: "hash",
			})
			require.NoError(t, err)
		}

		req1 := httptest.NewRequest(http.MethodGet, "/api/v1/history?limit=2", nil)
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		var res1 struct {
			Toggles    []map[string]any `json:"toggles"`
			NextCursor int64            `json:"next_cursor"`
			HasMore    bool             `json:"has_more"`
		}
		require.NoError(t, json.Unmarshal(rec1.Body.Bytes(), &res1))
		assert.Len(t, res1.Toggles, 2)
		assert.True(t, res1.HasMore)
		assert.True(t, res1.NextCursor > 0)

		req2 := httptest.NewRequest(http.MethodGet, "/api/v1/history?limit=2&before="+strconv.FormatInt(res1.NextCursor, 10), nil)
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusOK, rec2.Code)

		var res2 struct {
			Toggles    []map[string]any `json:"toggles"`
			NextCursor int64            `json:"next_cursor"`
			HasMore    bool             `json:"has_more"`
		}
		require.NoError(t, json.Unmarshal(rec2.Body.Bytes(), &res2))
		assert.Len(t, res2.Toggles, 1)
		assert.False(t, res2.HasMore)
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
