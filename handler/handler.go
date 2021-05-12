package handler

import (
	"fmt"
	"io"
	"log"
	"strings"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
	"github.com/ysugimoto/doorkeeper/github"
)

const (
	githubEventNamePullRequest = "pull_request"
	githubEventNamePing        = "ping"
	githubEventNamePush        = "push"

	githubPullRequestActionOpened      = "opened"
	githubPullRequestActionEdited      = "edited"
	githubPullRequestActionSynchronize = "synchronize"

	contextNameReleaseNote = "doorkeeper:releasenote"
	contextNameValidation  = "doorkeeper:validate"

	buildStatusPending = "pending"
	buildStatusFailure = "failure"
	buildStatusSuccess = "success"
)

func WebhookHandler(opts ...Option) http.Handler {
	c := &Config{}
	for _, o := range opts {
		o(c)
	}

	return http.StripPrefix(
		fmt.Sprintf("/%s", strings.Trim(c.prefix, "/")),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// Check webhook request comes from exact Github server
			if !compareSignature(r, c.appSecret) {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, "Signature unmatched")
				return
			}

			// Switch action by header
			switch r.Header.Get("X-Github-Event") {
			case githubEventNamePullRequest:
				handlePullRequestEvent(c, w, r)
			case githubEventNamePush:
				handlePushEvent(c, w, r)
			case githubEventNamePing:
				// Accept Ping event
				successResponse(w)
				return
			}

			// Forbid other events
			log.Println("We don't support event of '" + r.Header.Get("X-Github-Event") + "'")
			successResponse(w)
		}),
	)
}

func handlePullRequestEvent(c *Config, w http.ResponseWriter, r *http.Request) {
	// Accept PullRequest event
	var evt entity.GithubPullRequestEvent
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Failed to decode github webhook body to JSON "+err.Error())
		return
	}

	client := c.Client()

	// Get and parse rule from destination repository
	rr, err := client.RuleFile(r.Context(), evt.ContentURL(github.SettingFile), evt.Repository.FullName)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, ".doorkeeper.yml not found in "+evt.Repository.FullName)
		return
	}

	// switch actions by action
	switch evt.Action {

	// When new pullrequest has been opened, run validate and factory relates note
	case githubPullRequestActionOpened:
		if ok, _ := rr.MatchValidateBranch(evt.BaseBranch()); ok {
			if !rr.Validation.Disable {
				fmt.Println("execute validateion")
				go validatePullRequest(client, evt, rr)
			}
		}

		if ok, _ := rr.MatchReleaseNoteBranch(evt.BaseBranch()); ok {
			if !rr.ReleaseNote.Disable {
				fmt.Println("execute releasenote")
				go factoryRelaseNotes(client, evt, rr)
			}
		}

		// When pullrequest has been edited, only runs validate
	case githubPullRequestActionEdited:
		if ok, _ := rr.MatchValidateBranch(evt.BaseBranch()); ok {
			if !rr.Validation.Disable {
				fmt.Println("execute validate")
				go validatePullRequest(client, evt, rr)
			}
		}

		// When pullrequest has been synchronized, only runs factory release notes
	case githubPullRequestActionSynchronize:
		if ok, _ := rr.MatchReleaseNoteBranch(evt.BaseBranch()); ok {
			if !rr.ReleaseNote.Disable {
				fmt.Println("execute releasenote")
				go factoryRelaseNotes(client, evt, rr)
			}
		}
	}
	successResponse(w)
}

func handlePushEvent(c *Config, w http.ResponseWriter, r *http.Request) {
	// Accept push event
	var evt entity.GithubPushEvent
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Failed to decode github webhook body to JSON "+err.Error())
		return
	}

	client := c.Client()

	// Get and parse rule from destination repository
	rr, err := client.RuleFile(r.Context(), evt.ContentURL("/.doorkeeper.yml"), evt.Repository.FullName)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, ".doorkeeper.yml not found in "+evt.Repository.FullName)
		return
	}

	switch {
	case strings.HasPrefix(evt.Ref, "refs/tags"):
		if ok, _ := rr.MatchTag(strings.TrimPrefix(evt.Ref, "refs/tags/")); ok {
			if !rr.ReleaseNote.Disable {
				go processTagPushEvent(client, evt, rr)
			}
		}
	}
	successResponse(w)
}
