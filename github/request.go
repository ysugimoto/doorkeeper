package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"encoding/json"
	"net/http"
)

func sendGithubRequest(
	ctx context.Context,
	method string,
	url string,
	body interface{},
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

	log.Printf("Request API: %s\n", url)
	req, err := http.NewRequestWithContext(ctx, method, url, b)
	if err != nil {
		log.Printf("Request error: %s\n", err)
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
	req.Header.Set("Content-Type", "application/json")
	if acceptHeader != "" {
		req.Header.Set("Accept", acceptHeader)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		b := new(bytes.Buffer)
		io.Copy(b, resp.Body)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, b.String())
	}
	return resp, nil
}
