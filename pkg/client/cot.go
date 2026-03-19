package client

import (
	"context"
	"fmt"
	"strings"
)

// PublishCoT publishes a CoT XML event to the server.
func (c *Client) PublishCoT(ctx context.Context, eventXML []byte) error {
	if len(strings.TrimSpace(string(eventXML))) == 0 {
		return fmt.Errorf("CoT event XML is empty")
	}

	req, err := c.newRequest(ctx, "POST", c.cotEventsPath, "application/xml", eventXML)
	if err != nil {
		return err
	}

	_, err = c.do(req)
	return err
}
