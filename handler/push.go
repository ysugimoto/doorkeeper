package handler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
	"github.com/ysugimoto/doorkeeper/github"
	"github.com/ysugimoto/doorkeeper/rule"
)

const (
	integrationTypeSlack = "slack"
)

// goroutine
func processTagPushEvent(c *github.Client, evt entity.GithubPushEvent, r *rule.Rule) {
	ctx, timeout := context.WithTimeout(context.Background(), 3*time.Minute)
	defer timeout()

	currentTag := evt.CurrentTag()
	if currentTag == nil {
		log.Printf("Failed to parse current tag from ref: %s", evt.Ref)
		return
	}

	tags, err := c.Tags(ctx, evt.TagsURL())
	if err != nil {
		log.Printf("Failed to retrieve tag lists: %s", err)
		return
	}

	previousTag := tags.PreviousTag(currentTag)
	if previousTag == nil {
		previousTag = &entity.Tag{
			Raw: evt.Repository.MasterBranch,
		}
	}

	commits, err := c.Compare(ctx, evt.CompareURL(currentTag, previousTag))
	if err != nil {
		log.Printf("Failed to compare refs: %s, error: %s\n", evt.CompareURL(currentTag, previousTag), err)
		return
	}

	// Factory related notes from merged PullRequest description
	// Note that commit may duplicate so we need to be unique all commits by checking commit sha.
	var notes []string
	stack := make(map[int]struct{})
	for i := range commits {
		sha := commits[i].Sha
		prs, err := c.PullRequests(ctx, evt.PullRequestURL(sha))
		if err != nil {
			log.Printf("Failed to get commit-related pullrequests: %s, error: %s\n", evt.PullRequestURL(sha), err)
			return
		}
		for j := range prs {
			// Guard duplication
			if _, ok := stack[prs[j].Number]; ok {
				continue
			}
			stack[prs[j].Number] = struct{}{}
			matches := releaseSectionRegex.FindStringSubmatch(prs[j].Body)
			if matches != nil {
				notes = append(notes, fmt.Sprintf("- #%d %s", prs[j].Number, strings.TrimSpace(matches[1])))
			}
		}
	}

	message := fmt.Sprintf(":robot: There are no release notes found for tag: %s", currentTag.Raw)
	if len(notes) > 0 {
		message = fmt.Sprintf(
			":robot: Release Note collected succesfully for tag %s:\n\n```\n%s\n```\n\nTake this!",
			currentTag.Raw,
			strings.Join(notes, "\n"),
		)
	}

	integration := r.Integration()
	for k, v := range integration {
		switch k {
		case integrationTypeSlack:
			if err := sendToSlack(ctx, v, message); err != nil {
				log.Printf("Failed to send slack notification, error: %s\n", err)
				return
			}
		}
	}
}

func sendToSlack(ctx context.Context, webhookURL, message string) error {
	body, err := json.Marshal(map[string]string{
		"text": message,
	})
	if err != nil {
		return fmt.Errorf("Failed to marshal body: %w", err)
	}

	c, timeout := context.WithTimeout(ctx, 5*time.Second)
	defer timeout()

	req, err := http.NewRequestWithContext(c, http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("Failed to make slack request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to get slack response: %w", err)
	}
	resp.Body.Close()
	return nil
}
