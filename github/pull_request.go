package github

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

func PullRequests(ctx context.Context, url string) ([]entity.GithubPullRequest, error) {
	resp, err := sendGithubRequest(
		ctx,
		http.MethodGet,
		url,
		nil,
		"application/vnd.github.groot-preview+json",
	)

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
