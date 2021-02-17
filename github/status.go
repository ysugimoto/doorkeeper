package github

import (
	"context"
	"fmt"

	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

func Status(ctx context.Context, url string, s entity.GithubStatus) error {
	resp, err := sendGithubRequest(
		ctx,
		http.MethodPost,
		url,
		s,
		"application/vnd.github.v3+json",
	)

	if err != nil {
		return fmt.Errorf("Failed to call status request: %w", err)
	}
	resp.Body.Close()
	return nil
}
