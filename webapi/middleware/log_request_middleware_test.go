package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestContext(t *testing.T) *gin.Context {
	t.Helper()

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)
	if err := engine.SetTrustedProxies(nil); err != nil {
		t.Fatalf("set trusted proxies: %v", err)
	}
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return ctx
}

func TestReadClientIPUsesXRealForPrivateRemoteIP(t *testing.T) {
	ctx := newTestContext(t)
	ctx.Request.RemoteAddr = "10.1.2.3:1234"
	ctx.Request.Header.Set("X-Real-IP", "203.0.113.10")

	got := readClientIP(ctx)
	if got != "203.0.113.10" {
		t.Fatalf("readClientIP: expected X-Real-IP, got %q", got)
	}
}

func TestReadClientIPFallsBackWhenNoXRealIP(t *testing.T) {
	ctx := newTestContext(t)
	ctx.Request.RemoteAddr = "10.1.2.3:1234"

	got := readClientIP(ctx)
	if got != "10.1.2.3" {
		t.Fatalf("readClientIP: expected remote IP, got %q", got)
	}
}

func TestReadClientIPIgnoresXRealForPublicRemoteIP(t *testing.T) {
	ctx := newTestContext(t)
	ctx.Request.RemoteAddr = "8.8.8.8:443"
	ctx.Request.Header.Set("X-Real-IP", "203.0.113.10")

	got := readClientIP(ctx)
	if got != "8.8.8.8" {
		t.Fatalf("readClientIP: expected remote IP, got %q", got)
	}
}
