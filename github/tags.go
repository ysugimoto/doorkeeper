package github

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

func Tags(ctx context.Context, url string) (entity.Tags, error) {
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
	var body []struct {
		Ref    string `json:"ref"`
		Object struct {
			Sha string `json:"sha"`
		} `json:"object"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	tags := make(entity.Tags, 0, len(body))
	for i := len(body) - 1; i > 0; i-- {
		if tag := entity.ParseTag(body[i].Ref, body[i].Object.Sha); tag != nil {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}
