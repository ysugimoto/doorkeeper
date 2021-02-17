package github

import (
	"context"
	"fmt"

	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

func Review(ctx context.Context, url string, r entity.GithubReview) error {
	resp, err := sendGithubRequest(
		ctx,
		http.MethodPost,
		url,
		r,
		"application/vnd.github.v3+json",
	)

	if err != nil {
		return fmt.Errorf("Failed to call review request: %w", err)
	}
	resp.Body.Close()
	return nil
}
