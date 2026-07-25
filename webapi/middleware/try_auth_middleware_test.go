package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/webapi/auth"
	"github.com/xxxsen/common/webapi/proxyutil"
)

type testAuthenticator struct {
	name  string
	user  string
	err   error
	calls *[]string
}

func (a *testAuthenticator) Name() string {
	return a.name
}

func (a *testAuthenticator) Auth(
	_ *gin.Context,
	_ auth.UserQueryFunc,
) (string, error) {
	*a.calls = append(*a.calls, a.name)
	return a.user, a.err
}

func TestTryAuthMiddlewareUsesExplicitOrderAndStopsOnSuccess(t *testing.T) {
	var calls []string
	first := &testAuthenticator{
		name:  "first",
		err:   errors.New("not matched"),
		calls: &calls,
	}
	second := &testAuthenticator{name: "second", user: "alice", calls: &calls}
	third := &testAuthenticator{name: "third", user: "bob", calls: &calls}
	engine := gin.New()
	engine.Use(TryAuthMiddleware(noUsers, first, second, third))
	engine.GET("/", func(ctx *gin.Context) {
		info, ok := proxyutil.GetUserInfo(ctx.Request.Context())
		if !ok {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		ctx.String(http.StatusOK, info.AuthType+":"+info.Username)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/",
		nil,
	)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || recorder.Body.String() != "second:alice" {
		t.Fatalf("response = %d %q", recorder.Code, recorder.Body.String())
	}
	if len(calls) != 2 || calls[0] != "first" || calls[1] != "second" {
		t.Fatalf("calls = %#v", calls)
	}
}

func TestTryAuthMiddlewareDefaultsToBasic(t *testing.T) {
	engine := gin.New()
	engine.Use(TryAuthMiddleware(auth.MapUserMatch(map[string]string{
		"alice": "secret",
	})))
	engine.GET("/", func(ctx *gin.Context) {
		info, ok := proxyutil.GetUserInfo(ctx.Request.Context())
		if !ok {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		ctx.String(http.StatusOK, info.AuthType+":"+info.Username)
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/",
		nil,
	)
	request.SetBasicAuth("alice", "secret")
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || recorder.Body.String() != "basic:alice" {
		t.Fatalf("response = %d %q", recorder.Code, recorder.Body.String())
	}
}

func noUsers(context.Context, string) (string, bool, error) {
	return "", false, nil
}
