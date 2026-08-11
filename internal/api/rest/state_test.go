package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
			State bool `json:"state"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.False(t, res.State)
	})
}
