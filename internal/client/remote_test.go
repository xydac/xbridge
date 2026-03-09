package client

import "testing"

func TestResolveRemoteDefault(t *testing.T) {
	cfg := &ClientConfig{
		DefaultRemote: "mac:7900",
		Key:           "secret",
		Remotes:       map[string]Remote{},
	}

	host, key := cfg.ResolveRemote("")
	if host != "mac:7900" {
		t.Errorf("expected mac:7900, got %s", host)
	}
	if key != "secret" {
		t.Errorf("expected secret, got %s", key)
	}
}

func TestResolveRemoteNamed(t *testing.T) {
	cfg := &ClientConfig{
		DefaultRemote: "mac:7900",
		Key:           "default-key",
		Remotes: map[string]Remote{
			"admin": {Host: "mac:7901", Key: "admin-key"},
		},
	}

	host, key := cfg.ResolveRemote("admin")
	if host != "mac:7901" {
		t.Errorf("expected mac:7901, got %s", host)
	}
	if key != "admin-key" {
		t.Errorf("expected admin-key, got %s", key)
	}
}

func TestResolveRemoteNamedFallbackKey(t *testing.T) {
	cfg := &ClientConfig{
		DefaultRemote: "mac:7900",
		Key:           "default-key",
		Remotes: map[string]Remote{
			"main": {Host: "mac:7900"},
		},
	}

	host, key := cfg.ResolveRemote("main")
	if host != "mac:7900" {
		t.Errorf("expected mac:7900, got %s", host)
	}
	if key != "default-key" {
		t.Errorf("expected default-key, got %s", key)
	}
}

func TestResolveRemoteNotFound(t *testing.T) {
	cfg := &ClientConfig{
		DefaultRemote: "mac:7900",
		Key:           "key",
		Remotes:       map[string]Remote{},
	}

	host, key := cfg.ResolveRemote("nonexistent")
	if host != "mac:7900" {
		t.Errorf("expected fallback to default, got %s", host)
	}
	if key != "key" {
		t.Errorf("expected fallback key, got %s", key)
	}
}
