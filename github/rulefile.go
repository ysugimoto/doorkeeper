package github

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/rule"
	"gopkg.in/yaml.v2"
)

// RuleFile get content file from repostory and make rule
func (c *Client) RuleFile(ctx context.Context, u string) (*rule.Rule, error) {
	resp, err := c.apiRequest(ctx, http.MethodGet, u, nil, githubBasicHeader)
	if err != nil {
		return nil, fmt.Errorf("Failed to call content request: %s,  %w", u, err)
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
