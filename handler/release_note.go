package handler

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/ysugimoto/doorkeeper/entity"
	"github.com/ysugimoto/doorkeeper/github"
)

var releaseSectionRegex = regexp.MustCompile(`(?s)<!--\s*RELEASE\s*-->(.+?)<!--\s*/RELEASE\s*-->`)

// goroutine
func factoryRelaseNotes(c *github.Client, evt entity.GithubPullRequestEvent) {
	ctx, timeout := context.WithTimeout(context.Background(), 10*time.Minute)
	defer timeout()

	// Firstly, create status as "pending"
	err := c.Status(ctx, evt.StatusURL(), entity.GithubStatus{
		Status:      "pending",
		Context:     "grc:relasenote",
		Description: "factoty release note",
	})
	if err != nil {
		log.Println("Failed to create status as pending:", err)
		return
	}

	var factoryErr error
	defer func() {
		if factoryErr != nil {
			// Update to "failure" status
			c.Status(ctx, evt.StatusURL(), entity.GithubStatus{
				Status:      "failure",
				Context:     "grc:relasenote",
				Description: "factory release note",
			})
			return
		}
		// Otherwise, update to "success"
		c.Status(ctx, evt.StatusURL(), entity.GithubStatus{
			Status:      "success",
			Context:     "grc:relasenote",
			Description: "factory release note",
		})
	}()

	// Step1. compare refs in order to get included PRs
	commits, err := c.Compare(ctx, evt.CompareURL())
	if err != nil {
		log.Printf("Failed to compare refs: %s, error: %s\n", evt.CompareURL(), err)
		factoryErr = err
		return
	}

	// Step2. factory related notes from merged PullRequest description
	// Note that commit may duplicate so we need to be unique all commits by checking commit sha.
	var notes []string
	stack := make(map[int]struct{})
	for i := range commits {
		sha := commits[i].Sha
		prs, err := c.PullRequests(ctx, evt.PullRequestURL(sha))
		if err != nil {
			log.Printf("Failed to get commit-related pullrequests: %s, error: %s\n", evt.PullRequestURL(sha), err)
			factoryErr = err
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

	message := ":robot: There are no release notes found."
	if len(notes) > 0 {
		message = fmt.Sprintf(
			":robot: Release Note collected succesfully:\n\n```\n%s\n```\n\nTake this!",
			strings.Join(notes, "\n"),
		)
	}

	// And add review comment with release note
	c.Review(ctx, evt.ReviewURL(), entity.GithubReview{
		Body:  message,
		Event: "COMMENT",
	})
}
