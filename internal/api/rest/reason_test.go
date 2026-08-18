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

	t.Run("rejects reason exceeding 60 characters with 400", func(t *testing.T) {
		env := setupTestEnv(t)

		validUUID := uuid.NewString()
		longReason := strings.Repeat("a", 61)
		body := bytes.NewBufferString(`{"id":"` + validUUID + `","reason":"` + longReason + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("rejects reason containing URLs with 400", func(t *testing.T) {
		env := setupTestEnv(t)
		validUUID := uuid.NewString()

		testCases := []string{
			"check out https://example.com for info",
			"visit http://spam.net",
			"go to www.google.com please",
			"check google.com for details",
			"visit my-app.io now",
			"look at crypto.xyz",
			"join discord.gg/server",
		}

		for _, tc := range testCases {
			body := bytes.NewBufferString(`{"id":"` + validUUID + `","reason":"` + tc + `"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			env.router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "urls are not allowed")
		}
	})

	t.Run("rejects profane reasons with 400", func(t *testing.T) {
		env := setupTestEnv(t)
		validUUID := uuid.NewString()

		testCases := []string{
			"fuck this lamp",
			"sh!t happens",
			"you are an @$$hole",
		}

		for _, tc := range testCases {
			body := bytes.NewBufferString(`{"id":"` + validUUID + `","reason":"` + tc + `"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			env.router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "profanity is not allowed")
		}
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

	t.Run("rejects whitespace-only reason with 400", func(t *testing.T) {
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

		body := bytes.NewBufferString(`{"id":"` + toggleRes.ID + `","reason":"    "}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("rejects secondary reason submission attempt with 400 when reason is already set", func(t *testing.T) {
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

		body1 := bytes.NewBufferString(`{"id":"` + toggleRes.ID + `","reason":"first reason"}`)
		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body1)
		req1.Header.Set("Content-Type", "application/json")
		rec1 := httptest.NewRecorder()
		env.router.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		body2 := bytes.NewBufferString(`{"id":"` + toggleRes.ID + `","reason":"second reason"}`)
		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body2)
		req2.Header.Set("Content-Type", "application/json")
		rec2 := httptest.NewRecorder()
		env.router.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusBadRequest, rec2.Code)
	})

	t.Run("returns 500 when database failure occurs during reason update", func(t *testing.T) {
		env := setupTestEnv(t)
		validUUID := uuid.NewString()

		_ = env.db.Close() // Close database to simulate internal DB failure

		body := bytes.NewBufferString(`{"id":"` + validUUID + `","reason":"reason text"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reason", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		env.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
