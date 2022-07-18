package traefik_header_whitelist_plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	traefik_header_whitelist_plugin "github.com/icanbwell/cie.traefik-header-whitelist-plugin"
)

func TestAccept(t *testing.T) {
	ctx := context.Background()
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { nextCalled = true })
	cfg := traefik_header_whitelist_plugin.CreateConfig()
	cfg.Rules = []traefik_header_whitelist_plugin.Rule{
		{
			Key: "header-a",
			Value: "value-a",
		},
		{
			Key: "header-b",
			Value: "value-b",
		},
	}

	plugin, err := traefik_header_whitelist_plugin.New(ctx, next, cfg, "test-traefik-header-whitelist-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("header-a", "value-a")
	req.Header.Set("header-b", "value-b")

	plugin.ServeHTTP(recorder, req)

	if recorder.Code == http.StatusForbidden {
		t.Fatal("Expected OK")
	}

	if nextCalled == false {
		t.Fatal("next.ServeHTTP should be called, but wasn't")
	}
}

func TestDeny(t *testing.T) {
	ctx := context.Background()
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { nextCalled = true })
	cfg := traefik_header_whitelist_plugin.CreateConfig()
	cfg.Rules = []traefik_header_whitelist_plugin.Rule{
		{
			Key: "header-a",
			Value: "value-a",
		},
	}

	plugin, err := traefik_header_whitelist_plugin.New(ctx, next, cfg, "test-traefik-header-whitelist-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	plugin.ServeHTTP(recorder, req)

	if recorder.Code == http.StatusOK {
		t.Fatal("Expected Forbidden")
	}

	if nextCalled == true {
		t.Fatal("next.ServeHTTP should not be called, but was")
	}

}
