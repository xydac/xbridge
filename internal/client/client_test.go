package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/health" {
			t.Errorf("expected /health, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	c := New(server.URL, "")
	resp, err := c.Get("/health")
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]string
	if err := ReadJSON(resp, &result); err != nil {
		t.Fatal(err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected ok, got %s", result["status"])
	}
}

func TestClientPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("expected X-API-Key test-key, got %s", r.Header.Get("X-API-Key"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": "build-123"})
	}))
	defer server.Close()

	c := New(server.URL, "test-key")
	resp, err := c.Post("/build", map[string]string{"profile": "staging"})
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]string
	if err := ReadJSON(resp, &result); err != nil {
		t.Fatal(err)
	}
	if result["id"] != "build-123" {
		t.Errorf("expected build-123, got %s", result["id"])
	}
}

func TestClientReadJSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	c := New(server.URL, "")
	resp, err := c.Get("/missing")
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]string
	err = ReadJSON(resp, &result)
	if err == nil {
		t.Error("expected error for 404 response")
	}
}

func TestNewClientAddsHTTP(t *testing.T) {
	c := New("mac:7900", "key")
	if c.BaseURL != "http://mac:7900" {
		t.Errorf("expected http:// prefix, got %s", c.BaseURL)
	}
}

func TestNewClientKeepsHTTPS(t *testing.T) {
	c := New("https://mac:7900", "key")
	if c.BaseURL != "https://mac:7900" {
		t.Errorf("expected https preserved, got %s", c.BaseURL)
	}
}
