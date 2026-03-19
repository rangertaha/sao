package client

import (
	"fmt"
	"net/http"
	"strings"
)

// Option configures a client instance.
type Option func(*Client) error

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(httpClient HTTPDoer) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("HTTP client cannot be nil")
		}

		c.httpClient = httpClient
		return nil
	}
}

// WithBearerToken configures bearer auth for all requests.
func WithBearerToken(token string) Option {
	return func(c *Client) error {
		c.bearerToken = strings.TrimSpace(token)
		return nil
	}
}

// WithUserAgent configures the user agent header.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) error {
		if strings.TrimSpace(userAgent) == "" {
			return fmt.Errorf("user agent cannot be empty")
		}

		c.userAgent = userAgent
		return nil
	}
}

// WithHealthPath overrides the health endpoint path.
func WithHealthPath(path string) Option {
	return func(c *Client) error {
		path = strings.TrimSpace(path)
		if path == "" || path[0] != '/' {
			return fmt.Errorf("health path must start with '/'")
		}

		c.healthPath = path
		return nil
	}
}

// WithCoTEventsPath overrides the CoT events endpoint path.
func WithCoTEventsPath(path string) Option {
	return func(c *Client) error {
		path = strings.TrimSpace(path)
		if path == "" || path[0] != '/' {
			return fmt.Errorf("CoT events path must start with '/'")
		}

		c.cotEventsPath = path
		return nil
	}
}

// Compile-time guard for common use.
var _ HTTPDoer = (*http.Client)(nil)
