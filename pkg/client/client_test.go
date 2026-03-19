package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/v1/health" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	c, err := New(server.URL)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	status, err := c.Health(context.Background())
	if err != nil {
		t.Fatalf("Health() error: %v", err)
	}

	if status.Status != "ok" {
		t.Fatalf("unexpected health status: %s", status.Status)
	}
}

func TestReady(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/v1/ready" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ready":true,"checks":{"config_loaded":true,"nats_ready":true,"runtime_ready":true}}`))
	}))
	defer server.Close()

	c, err := New(server.URL)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	status, err := c.Ready(context.Background())
	if err != nil {
		t.Fatalf("Ready() error: %v", err)
	}
	if !status.Ready {
		t.Fatal("expected ready=true")
	}
	if !status.Checks["config_loaded"] {
		t.Fatal("expected config_loaded check true")
	}
}

func TestPublishCoT(t *testing.T) {
	t.Parallel()

	const expectedToken = "token-123"
	const event = `<event version="2.0" uid="A"/>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/v1/cot/events" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+expectedToken {
			t.Fatalf("unexpected auth header: %s", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/xml" {
			t.Fatalf("unexpected content-type: %s", got)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed reading body: %v", err)
		}
		if string(body) != event {
			t.Fatalf("unexpected body: %s", string(body))
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	c, err := New(server.URL, WithBearerToken(expectedToken))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	if err := c.PublishCoT(context.Background(), []byte(event)); err != nil {
		t.Fatalf("PublishCoT() error: %v", err)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid event"))
	}))
	defer server.Close()

	c, err := New(server.URL)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	err = c.PublishCoT(context.Background(), []byte(`<event/>`))
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError but got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", apiErr.StatusCode)
	}
}
