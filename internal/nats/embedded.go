package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/rangertaha/sao/internal/config"
)

// Embedded wraps an in-process NATS server.
type Embedded struct {
	server *server.Server
}

// Start boots an embedded NATS server and binds its lifecycle to ctx.
func Start(ctx context.Context, cfg config.NATSConfig) (*Embedded, error) {
	opts := &server.Options{
		Host:   cfg.Host,
		Port:   cfg.Port,
		NoSigs: true,
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("create nats server: %w", err)
	}

	go ns.Start()
	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, fmt.Errorf("nats server not ready on %s", ns.ClientURL())
	}

	runtime := &Embedded{server: ns}

	go func() {
		<-ctx.Done()
		ns.Shutdown()
	}()

	return runtime, nil
}

// URL returns the client URL for this embedded NATS server.
func (e *Embedded) URL() string {
	if e == nil || e.server == nil {
		return ""
	}
	return e.server.ClientURL()
}

// Shutdown gracefully terminates embedded NATS.
func (e *Embedded) Shutdown(ctx context.Context) error {
	if e == nil || e.server == nil {
		return nil
	}

	e.server.Shutdown()

	done := make(chan struct{})
	go func() {
		e.server.WaitForShutdown()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
