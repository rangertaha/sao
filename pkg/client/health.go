package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// HealthStatus is returned by the health endpoint.
type HealthStatus struct {
	Status string `json:"status"`
}

// ReadyStatus is returned by the readiness endpoint.
type ReadyStatus struct {
	Ready  bool            `json:"ready"`
	Checks map[string]bool `json:"checks"`
}

// Health requests server health.
func (c *Client) Health(ctx context.Context) (*HealthStatus, error) {
	req, err := c.newRequest(ctx, "GET", c.healthPath, "", nil)
	if err != nil {
		return nil, err
	}

	payload, err := c.do(req)
	if err != nil {
		return nil, err
	}

	var out HealthStatus
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil, fmt.Errorf("decode health response: %w", err)
	}

	return &out, nil
}

// Ready requests server readiness.
func (c *Client) Ready(ctx context.Context) (*ReadyStatus, error) {
	req, err := c.newRequest(ctx, "GET", "/v1/ready", "", nil)
	if err != nil {
		return nil, err
	}

	payload, err := c.do(req)
	if err != nil {
		return nil, err
	}

	var out ReadyStatus
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil, fmt.Errorf("decode ready response: %w", err)
	}

	return &out, nil
}
