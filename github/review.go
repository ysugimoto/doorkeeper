package github

import (
	"context"
	"fmt"

	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

// Review create review comment to destination pullrequest
func (c *Client) Review(ctx context.Context, url, repository string, r entity.GithubReview) error {
	resp, err := c.apiRequest(ctx, http.MethodPost, url, r, repository, githubBasicHeader)
	if err != nil {
		return fmt.Errorf("Failed to call review request: %w", err)
	}

	resp.Body.Close()
	return nil
}
