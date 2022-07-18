package traefik_header_whitelist_plugin

import (
	"context"
	"net/http"
)

// Header key and values to allow
type Rule struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Config the plugin configuration.
type Config struct {
	Rules []Rule `json:"rules"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// Example a plugin.
type HeaderWhitelist struct {
	next  http.Handler
	rules []Rule
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// ...
	return &HeaderWhitelist{
		next:  next,
		rules: config.Rules,
	}, nil
}

func (e *HeaderWhitelist) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	failed := false
	for i := 0; i < len(e.rules); i += 1 {
		v := req.Header.Get(e.rules[i].Key)
		if v != e.rules[i].Value {
			failed = true
			rw.WriteHeader(http.StatusForbidden)
			break
		}
	}
	if !failed {
		rw.WriteHeader(http.StatusOK)
        e.next.ServeHTTP(rw, req)
	}
}
