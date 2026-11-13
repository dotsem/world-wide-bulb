package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"world-wide-bulb/internal/api/middleware"

	"github.com/stretchr/testify/assert"
)

func TestTurnstileVerifier(t *testing.T) {
	fake := func(t *testing.T, status int, body string) string {
		t.Helper()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			assert.NoError(t, r.ParseForm())
			assert.Equal(t, "the-secret", r.PostForm.Get("secret"))
			assert.Equal(t, "the-token", r.PostForm.Get("response"))
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		}))
		t.Cleanup(srv.Close)
		return srv.URL
	}

	cases := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{"accepts a successful verification", http.StatusOK, `{"success":true}`, true},
		{"rejects a failed verification", http.StatusOK, `{"success":false,"error-codes":["invalid-input-response"]}`, false},
		{"rejects a body without success", http.StatusOK, `{}`, false},
		{"rejects malformed JSON", http.StatusOK, `not json`, false},
		{"rejects a server error", http.StatusInternalServerError, `oops`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			verify := middleware.TurnstileVerifier("the-secret", fake(t, tc.status, tc.body))
			assert.Equal(t, tc.want, verify(context.Background(), "the-token"))
		})
	}

	t.Run("rejects when siteverify is unreachable", func(t *testing.T) {
		srv := httptest.NewServer(http.NotFoundHandler())
		srv.Close()
		verify := middleware.TurnstileVerifier("the-secret", srv.URL)
		assert.False(t, verify(context.Background(), "the-token"))
	})

	t.Run("rejects an invalid verify URL", func(t *testing.T) {
		verify := middleware.TurnstileVerifier("the-secret", "://bad")
		assert.False(t, verify(context.Background(), "the-token"))
	})

	t.Run("rejects when the request context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		verify := middleware.TurnstileVerifier("the-secret", fake(t, http.StatusOK, `{"success":true}`))
		assert.False(t, verify(ctx, "the-token"))
	})
}
