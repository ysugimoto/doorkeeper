package github

import (
	"context"
	"fmt"

	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/rule"
	"gopkg.in/yaml.v2"
)

// RuleFile get content file from repostory and make rule
func (c *Client) RuleFile(ctx context.Context, url, repostory string) (*rule.Rule, error) {
	resp, err := c.apiRequest(ctx, http.MethodGet, url, nil, repostory, githubBasicHeader)
	if err != nil {
		return nil, fmt.Errorf("Failed to call content request: %s,  %w", url, err)
	}
	defer resp.Body.Close()

	var content struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}
	data, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode content: %w", err)
	}

	var r rule.Rule
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal yaml: %w", err)
	}

	return &r, nil
}
