package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultHealthPath = "/v1/health"
	defaultCoTPath    = "/v1/cot/events"
	defaultUserAgent  = "sao-client/0.1"
)

// HTTPDoer is implemented by *http.Client and test doubles.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client provides a reusable SAO server API client.
type Client struct {
	baseURL       *url.URL
	httpClient    HTTPDoer
	bearerToken   string
	userAgent     string
	healthPath    string
	cotEventsPath string
}

// New creates a new SAO client.
func New(baseURL string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("base URL must include scheme and host")
	}

	c := &Client{
		baseURL:       parsedURL,
		httpClient:    http.DefaultClient,
		userAgent:     defaultUserAgent,
		healthPath:    defaultHealthPath,
		cotEventsPath: defaultCoTPath,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, method, path, contentType string, body []byte) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	targetURL, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("resolve request path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	}

	return req, nil
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       payload,
		}
	}

	return payload, nil
}

// APIError wraps a non-2xx HTTP response.
type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("server returned status %d: %s", e.StatusCode, strings.TrimSpace(string(e.Body)))
}
