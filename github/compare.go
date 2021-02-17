package github

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

func Compare(ctx context.Context, url string) ([]entity.GithubCommit, error) {
	resp, err := sendGithubRequest(
		ctx,
		http.MethodGet,
		url,
		nil,
		"application/vnd.github.v3+json",
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to call compare request: %w", err)
	}
	defer resp.Body.Close()
	var commits struct {
		Commits []entity.GithubCommit `json:"commits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return commits.Commits, nil
}
