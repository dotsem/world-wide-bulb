package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostReason(t *testing.T) {
	t.Run("successfully submits reason for valid toggle UUID", func(t *testing.T) {
		env := setupTestEnv(t)

		reqToggle := httptest.NewRequest(http.MethodPost, "/api/v1/toggle", nil)
		recToggle := httptest.NewRecorder()
		env.router.ServeHTTP(recToggle, reqToggle)
		require.Equal(t, http.StatusOK, recToggle.Code)

		var toggleRes struct {
			ID string `json:"id"`
		}
		err := json.Unmarshal(recToggle.Body.Bytes(), &toggleRes)
		require.NoError(t, err)
		require.NotEmpty(t, toggleRes.ID)

		body := bytes.NewBufferString(`{"id":"` + toggleRes.ID + `","reason":"felt like turning it on"}`)
		reqReason := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
		reqReason.Header.Set("Content-Type", "application/json")
		recReason := httptest.NewRecorder()
		env.router.ServeHTTP(recReason, reqReason)

		assert.Equal(t, http.StatusOK, recReason.Code)
	})

	t.Run("rejects invalid UUID format with 400", func(t *testing.T) {
		env := setupTestEnv(t)

		body := bytes.NewBufferString(`{"id":"not-a-uuid","reason":"valid reason"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("rejects reason exceeding 100 characters with 400", func(t *testing.T) {
		env := setupTestEnv(t)

		validUUID := uuid.NewString()
		longReason := strings.Repeat("a", 101)
		body := bytes.NewBufferString(`{"id":"` + validUUID + `","reason":"` + longReason + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 404 for non-existent toggle UUID", func(t *testing.T) {
		env := setupTestEnv(t)

		randomUUID := uuid.NewString()
		body := bytes.NewBufferString(`{"id":"` + randomUUID + `","reason":"some reason"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
