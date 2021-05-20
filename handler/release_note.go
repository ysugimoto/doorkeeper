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
	"github.com/ysugimoto/doorkeeper/rule"
)

var releaseSectionRegex = regexp.MustCompile(`(?s)<!--\s*RELEASE\s*-->(.+?)<!--\s*/RELEASE\s*-->`)

// goroutine
func factoryRelaseNotes(c *github.Client, evt entity.GithubPullRequestEvent, r *rule.Rule) {
	ctx, timeout := context.WithTimeout(context.Background(), 10*time.Minute)
	defer timeout()

	// Firstly, create status as "pending"
	err := c.Status(ctx, evt.StatusURL(), evt.Repository.FullName, entity.GithubStatus{
		Status:      buildStatusPending,
		Context:     contextNameReleaseNote,
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
			c.Status(ctx, evt.StatusURL(), evt.Repository.FullName, entity.GithubStatus{
				Status:      buildStatusFailure,
				Context:     contextNameReleaseNote,
				Description: "factory release note",
			})
			return
		}
		// Otherwise, update to "success"
		c.Status(ctx, evt.StatusURL(), evt.Repository.FullName, entity.GithubStatus{
			Status:      buildStatusSuccess,
			Context:     contextNameReleaseNote,
			Description: "factory release note",
		})
	}()

	// Step1. compare refs in order to get included PRs
	commits, err := c.Compare(ctx, evt.CompareURL(), evt.Repository.FullName)
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
		prs, err := c.PullRequests(ctx, evt.PullRequestURL(sha), evt.Repository.FullName)
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
				notes = append(notes, fmt.Sprintf(
					"- #%d %s",
					prs[j].Number,
					formatReleaseNoteText(matches[1]),
				))
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

	// And add review comment with release note
	c.Review(ctx, evt.ReviewURL(), evt.Repository.FullName, entity.GithubReview{
		Body:  message,
		Event: "COMMENT",
	})
}
