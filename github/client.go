package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httputil"
)

const (
	githubBasicHeader       = "application/vnd.github.v3+json"
	githubRootPreviewHeader = "application/vnd.github.groot-preview+json"

	SettingFile = "/.doorkeeper.yml"
)

type Client struct {
	client        *http.Client
	authorization Authorization
}

func NewClient(c *http.Client, options ...Option) *Client {
	gc := &Client{
		client: c,
	}

	for i := range options {
		options[i](gc)
	}

	return gc
}

// Call API request to Github
func (c *Client) apiRequest(
	ctx context.Context,
	method string,
	url string,
	body interface{},
	repository string,
	acceptHeader string,
) (*http.Response, error) {

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

	if c.authorization != nil {
		c.authorization.Header(req, repository)
	}

	req.Header.Set("Content-Type", "application/json")
	if acceptHeader != "" {
		req.Header.Set("Accept", acceptHeader)
	}

	d, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(d))

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

func DefaultClient(options ...Option) *Client {
	return NewClient(http.DefaultClient, options...)
}
