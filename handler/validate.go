package handler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ysugimoto/doorkeeper/entity"
	"github.com/ysugimoto/doorkeeper/github"
	"github.com/ysugimoto/doorkeeper/rule"
)

// goroutine
func validatePullRequest(c *github.Client, evt entity.GithubPullRequestEvent, r *rule.Rule) {
	ctx, timeout := context.WithTimeout(context.Background(), 3*time.Minute)
	defer timeout()

	// Firstly, create status as "pending"
	if err := c.Status(ctx, evt.StatusURL(), evt.Repository.FullName, entity.GithubStatus{
		Status:      buildStatusPending,
		Context:     contextNameValidation,
		Description: "validate pull request",
	}); err != nil {
		log.Println("Failed to create status as pending:", err)
		return
	}

	var statusErr error
	defer func() {
		if statusErr != nil {
			// Update to "failure" status
			if err := c.Status(ctx, evt.StatusURL(), evt.Repository.FullName, entity.GithubStatus{
				Status:      buildStatusFailure,
				Context:     contextNameValidation,
				Description: "validate pull request",
			}); err != nil {
				log.Println("Failed to update check status as pending:", err)
			}
			// And add review comment what is invalid
			if err := c.Review(ctx, evt.ReviewURL(), evt.Repository.FullName, entity.GithubReview{
				Body:  statusErr.Error(),
				Event: "COMMENT",
			}); err != nil {
				log.Println("Failed to send comment:", err)
			}
			return
		}
		// Otherwise, update to "success"
		if err := c.Status(ctx, evt.StatusURL(), evt.Repository.FullName, entity.GithubStatus{
			Status:      buildStatusSuccess,
			Context:     contextNameValidation,
			Description: "validate pull request",
		}); err != nil {
			log.Println("Failed to update shcek status as success:", err)
		}
	}()

	errors := make([]string, 0, 2)
	if err := r.ValidateTitle(evt.PullRequest.Title); err != nil {
		errors = append(errors, "- "+err.Error())
	}
	if err := r.ValidateDescription(evt.PullRequest.Body); err != nil {
		errors = append(errors, "- "+err.Error())
	}

	if len(errors) > 0 {
		message := fmt.Sprintf(
			":robot: PR Validation Failed!\n%s\n",
			strings.Join(errors, "\n"),
		)
		statusErr = fmt.Errorf(message)

		// Also send to slack
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
	// passed

}
