package webapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/webapi/auth"
	"github.com/xxxsen/common/webapi/proxyutil"
)

type engineAuthenticator struct {
	name string
	user string
}

func (a engineAuthenticator) Name() string {
	return a.name
}

func (a engineAuthenticator) Auth(
	_ *gin.Context,
	_ auth.UserQueryFunc,
) (string, error) {
	return a.user, nil
}

func TestEngineAuthenticatorsAreIsolated(t *testing.T) {
	first := newAuthTestEngine(t, engineAuthenticator{name: "one", user: "alice"})
	second := newAuthTestEngine(t, engineAuthenticator{name: "two", user: "bob"})

	var wait sync.WaitGroup
	wait.Add(2)
	responses := make(chan string, 2)
	go serveAuthTestEngine(first, &wait, responses)
	go serveAuthTestEngine(second, &wait, responses)
	wait.Wait()
	close(responses)

	found := make(map[string]bool, 2)
	for response := range responses {
		found[response] = true
	}
	if !found["one:alice"] || !found["two:bob"] {
		t.Fatalf("responses = %#v", found)
	}
}

func TestEngineNoRouteHandlerIsEngineLocal(t *testing.T) {
	custom, err := NewEngine(
		"",
		"",
		WithNoRoute(func(ctx *gin.Context) {
			ctx.String(http.StatusTeapot, "custom")
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	standard, err := NewEngine("", "")
	if err != nil {
		t.Fatal(err)
	}

	customRecorder := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", nil)
	custom.ServeHTTP(customRecorder, request)
	if customRecorder.Code != http.StatusTeapot || customRecorder.Body.String() != "custom" {
		t.Fatalf("custom no-route response = %d %q", customRecorder.Code, customRecorder.Body.String())
	}

	standardRecorder := httptest.NewRecorder()
	request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", nil)
	standard.ServeHTTP(standardRecorder, request)
	if standardRecorder.Code != http.StatusNotFound {
		t.Fatalf("standard no-route status = %d", standardRecorder.Code)
	}
}

func newAuthTestEngine(t *testing.T, authenticator auth.IAuth) IWebEngine {
	t.Helper()
	engine, err := NewEngine(
		"",
		"",
		WithAuthenticators(authenticator),
		WithRegister(func(group *gin.RouterGroup) {
			group.GET("/", func(ctx *gin.Context) {
				info, ok := proxyutil.GetUserInfo(ctx.Request.Context())
				if !ok {
					ctx.Status(http.StatusUnauthorized)
					return
				}
				ctx.String(http.StatusOK, info.AuthType+":"+info.Username)
			})
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	return engine
}

func serveAuthTestEngine(
	engine IWebEngine,
	wait *sync.WaitGroup,
	responses chan<- string,
) {
	defer wait.Done()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/",
		nil,
	)
	engine.ServeHTTP(recorder, request)
	responses <- recorder.Body.String()
}
