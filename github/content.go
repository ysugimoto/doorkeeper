package github

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/rule"
	"gopkg.in/yaml.v2"
)

func Content(ctx context.Context, url string) (*rule.Rule, error) {
	resp, err := sendGithubRequest(
		ctx,
		http.MethodGet,
		url,
		nil,
		"application/vnd.github.v3+json",
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to call content request: %w", err)
	}
	defer resp.Body.Close()
	var content struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	var r rule.Rule
	if err := yaml.Unmarshal([]byte(content.Content), &r); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal yaml: %w", err)
	}

	return &r, nil
}
