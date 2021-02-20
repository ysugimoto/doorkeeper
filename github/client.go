package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"encoding/json"
	"net/http"
)

const (
	githubBasicHeader       = "application/vnd.github.v3+json"
	githubRootPreviewHeader = "application/vnd.github.groot-preview+json"

	SettingFile = "/.doorkeeper.yml"
)

type Client struct {
	Timeout time.Duration
	client  *http.Client
}

func NewClient(c *http.Client) *Client {
	return &Client{
		Timeout: 5 * time.Second,
		client:  c,
	}
}

// Call API request to Github
func (c *Client) apiRequest(
	ctx context.Context,
	method string,
	url string,
	body interface{},
	acceptHeader string,
) (*http.Response, error) {

	ctx, timeout := context.WithTimeout(ctx, c.Timeout)
	defer timeout()

	var b io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal request body: %w", err)
		}
		b = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, b)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to %s %s %w", method, url, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
	req.Header.Set("Content-Type", "application/json")
	if acceptHeader != "" {
		req.Header.Set("Accept", acceptHeader)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 400 {
		b := new(bytes.Buffer)
		io.Copy(b, resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, b.String())
	}

	return resp, nil
}

var DefaultClient = NewClient(http.DefaultClient)
