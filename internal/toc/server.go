package toc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rangertaha/sao/internal/config"
	embeddednats "github.com/rangertaha/sao/internal/nats"
)

// Server is the TOC server orchestrator.
type Server struct {
	cfg     *config.Config
	nats    *embeddednats.Embedded
	runtime *Runtime
}

// NewServer creates a TOC server.
func NewServer(cfg *config.Config, natsRuntime *embeddednats.Embedded) *Server {
	return &Server{
		cfg:     cfg,
		nats:    natsRuntime,
		runtime: NewRuntime(),
	}
}

// Run starts TOC server runtime.
func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/health", s.handleHealth)
	mux.HandleFunc("/v1/ready", s.handleReady)
	mux.HandleFunc("/v1/cot/events", s.handleCoTEvents)

	httpServer := &http.Server{
		Addr:    s.cfg.Server.Address,
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return httpServer.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := map[string]string{
		"status":   "ok",
		"nats_url": s.nats.URL(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checks := map[string]bool{
		"config_loaded": s.cfg != nil,
		"nats_ready":    s.nats != nil && s.nats.URL() != "",
		"runtime_ready": s.runtime != nil && s.runtime.Router != nil,
	}

	ready := true
	for _, ok := range checks {
		if !ok {
			ready = false
			break
		}
	}

	resp := map[string]any{
		"ready":  ready,
		"checks": checks,
	}

	w.Header().Set("Content-Type", "application/json")
	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleCoTEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const maxBodyBytes = 1 << 20 // 1MB safeguard.
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("read body: %v", err), http.StatusBadRequest)
		return
	}

	delivered, err := s.runtime.RouteRawEvent(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid cot event: %v", err), http.StatusBadRequest)
		return
	}

	resp := map[string]any{
		"status":    "accepted",
		"delivered": delivered,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(resp)
}
