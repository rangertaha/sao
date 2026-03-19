package ui

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// Serve starts serving embedded UI assets and stops on context cancellation.
func Serve(ctx context.Context, addr string) error {
	handler, err := AssetHandler()
	if err != nil {
		return fmt.Errorf("create asset handler: %w", err)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return server.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}
