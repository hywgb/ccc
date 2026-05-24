package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/divord97/ccc/pkg/response"
)

// Service-to-service auth header names.
const (
	HeaderServiceToken     = "X-CCC-Service-Token"
	HeaderServiceTimestamp = "X-CCC-Service-Timestamp"
	HeaderServiceTenantID  = "X-CCC-Tenant-ID"
)

// Allowed clock drift between caller and server when verifying timestamps.
const serviceAuthMaxSkew = 5 * time.Minute

// nonceCache prevents replay attacks within the validity window.
type nonceCache struct {
	mu      sync.Mutex
	entries map[string]time.Time
}

func newNonceCache() *nonceCache { return &nonceCache{entries: make(map[string]time.Time)} }

func (c *nonceCache) seen(token string, now time.Time) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, t := range c.entries {
		if now.Sub(t) > serviceAuthMaxSkew*2 {
			delete(c.entries, k)
		}
	}
	if _, ok := c.entries[token]; ok {
		return true
	}
	c.entries[token] = now
	return false
}

// ServiceAuth authenticates internal service-to-service calls (e.g. FreeSWITCH → API)
// via an HMAC-SHA256 signature over (timestamp + method + path + body).
//
// Required headers:
//
//	X-CCC-Service-Timestamp: unix seconds
//	X-CCC-Service-Token:     hex(HMAC-SHA256(secret, ts + "\n" + method + "\n" + path + "\n" + sha256(body)))
//	X-CCC-Tenant-ID:         tenant identifier propagated to context
//
// If secret is empty the middleware rejects all requests (fail-closed).
func ServiceAuth(secret string) func(http.Handler) http.Handler {
	cache := newNonceCache()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if secret == "" {
				response.Error(w, http.StatusServiceUnavailable, "service auth not configured")
				return
			}

			tsStr := r.Header.Get(HeaderServiceTimestamp)
			token := r.Header.Get(HeaderServiceToken)
			tenantStr := r.Header.Get(HeaderServiceTenantID)
			if tsStr == "" || token == "" {
				response.Error(w, http.StatusUnauthorized, "missing service auth headers")
				return
			}

			tsSec, err := strconv.ParseInt(tsStr, 10, 64)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid timestamp")
				return
			}
			now := time.Now()
			ts := time.Unix(tsSec, 0)
			if d := now.Sub(ts); d > serviceAuthMaxSkew || d < -serviceAuthMaxSkew {
				response.Error(w, http.StatusUnauthorized, "timestamp out of range")
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "cannot read request body")
				return
			}
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytesReader(body))
			r.ContentLength = int64(len(body))

			bodyHash := sha256.Sum256(body)
			payload := fmt.Sprintf("%d\n%s\n%s\n%s", tsSec, r.Method, r.URL.Path, hex.EncodeToString(bodyHash[:]))
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write([]byte(payload))
			expected := hex.EncodeToString(mac.Sum(nil))

			if !hmac.Equal([]byte(expected), []byte(token)) {
				response.Error(w, http.StatusUnauthorized, "invalid service token")
				return
			}

			if cache.seen(token, now) {
				response.Error(w, http.StatusUnauthorized, "replay detected")
				return
			}

			ctx := r.Context()
			if tenantStr != "" {
				if tenantID, err := strconv.ParseInt(tenantStr, 10, 64); err == nil {
					ctx = context.WithValue(ctx, ContextKeyTenantID, tenantID)
				}
			}
			ctx = context.WithValue(ctx, ContextKeyRole, "service")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
