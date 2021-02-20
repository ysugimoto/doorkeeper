package github

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

// PullRequests retreives related pullrequest from commit hash
func (c *Client) PullRequests(ctx context.Context, url string) ([]entity.GithubPullRequest, error) {
	resp, err := c.apiRequest(ctx, http.MethodGet, url, nil, githubRootPreviewHeader)
	if err != nil {
		return nil, fmt.Errorf("Failed to call pullrequests request: %w", err)
	}
	defer resp.Body.Close()

	var prs []entity.GithubPullRequest
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return prs, nil
}
