package toc

import (
	"context"
	"fmt"

	"github.com/rangertaha/sao/internal/cot"
)

// Runtime owns TOC runtime modules.
type Runtime struct {
	Router *cot.Router
}

// NewRuntime constructs runtime with modular CoT router.
func NewRuntime() *Runtime {
	subs := cot.NewSubscriptions()
	router := cot.NewRouter(subs)
	return &Runtime{Router: router}
}

// Run keeps runtime alive until context cancellation.
func (r *Runtime) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// RouteRawEvent parses and routes raw CoT XML payloads.
func (r *Runtime) RouteRawEvent(xmlPayload []byte) (int, error) {
	event, err := cot.ParseEvent(xmlPayload)
	if err != nil {
		return 0, fmt.Errorf("parse event: %w", err)
	}
	return r.Router.Route(event), nil
}
