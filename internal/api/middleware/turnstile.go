package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"world-wide-bulb/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	humanCookieName       = "human"
	humanIDKey            = "human_id"
	humanTTL              = time.Hour
	turnstileHeader       = "X-Turnstile-Token"
	turnstileMaxTokenSize = 2048
	turnstileVerifyURL    = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

// TurnstileVerifier returns a function that validates a Turnstile token against Cloudflare's siteverify API.
func TurnstileVerifier(secret string) func(ctx context.Context, token string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	return func(ctx context.Context, token string) bool {
		form := url.Values{"secret": {secret}, "response": {token}}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, turnstileVerifyURL, strings.NewReader(form.Encode()))
		if err != nil {
			return false
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res, err := client.Do(req)
		if err != nil {
			slog.Warn("turnstile siteverify unreachable", slog.Any("err", err))
			return false
		}
		defer func() { _ = res.Body.Close() }()
		var body struct {
			Success bool `json:"success"`
		}
		return json.NewDecoder(res.Body).Decode(&body) == nil && body.Success
	}
}

// HumanID returns the ID of the passed Turnstile challenge for this request, or "" if there is none.
// Each challenge gets its own ID, so clients sharing an IP are told apart without letting one
// client mint new identities for free.
func HumanID(c *gin.Context, hasher *utils.Hasher) string {
	if id := c.GetString(humanIDKey); id != "" {
		return id
	}
	raw, err := c.Cookie(humanCookieName)
	if err != nil {
		return ""
	}
	parts := strings.Split(raw, ".")
	if len(parts) != 3 || !hasher.Verify(humanPayload(c.ClientIP(), parts[0], parts[1]), parts[2]) {
		return ""
	}
	if exp, err := strconv.ParseInt(parts[1], 10, 64); err != nil || time.Now().Unix() >= exp {
		return ""
	}
	return parts[0]
}

// Human returns a middleware that rejects requests until the client has passed a Turnstile challenge.
// A passed challenge is remembered in a signed, IP-bound cookie for humanTTL.
//
// ponytail: the cookie can be copied into a script on the same IP until it expires; verify a
// fresh token on every toggle if that gets abused.
func Human(siteKey string, hasher *utils.Hasher, isProd bool, verify func(ctx context.Context, token string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if HumanID(c, hasher) != "" {
			c.Next()
			return
		}

		token := c.GetHeader(turnstileHeader)
		if token == "" || len(token) > turnstileMaxTokenSize || !verify(c.Request.Context(), token) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":              "human verification required",
				"turnstile_site_key": siteKey,
			})
			return
		}

		id := uuid.NewString()
		exp := strconv.FormatInt(time.Now().Add(humanTTL).Unix(), 10)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(humanCookieName, id+"."+exp+"."+hasher.Sign(humanPayload(c.ClientIP(), id, exp)), int(humanTTL.Seconds()), "/", "", isProd, true)
		c.Set(humanIDKey, id)
		c.Next()
	}
}

func humanPayload(ip, id, exp string) string {
	return "human|" + ip + "|" + id + "|" + exp
}
