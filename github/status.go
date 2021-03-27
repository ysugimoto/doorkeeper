package github

import (
	"context"
	"fmt"

	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

// Status update pullrequest status
func (c *Client) Status(ctx context.Context, url, repository string, s entity.GithubStatus) error {
	resp, err := c.apiRequest(ctx, http.MethodPost, url, s, repository, githubBasicHeader)
	if err != nil {
		return fmt.Errorf("Failed to call status request: %w", err)
	}

	resp.Body.Close()
	return nil
}
