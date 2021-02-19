package handler

import (
	"fmt"
	"io"
	"strings"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
	"github.com/ysugimoto/doorkeeper/github"
	"github.com/ysugimoto/doorkeeper/rule"
)

const (
	githubEventNamePullRequest = "pull_request"
	githubEventNamePing        = "ping"
	githubEventNamePush        = "push"

	githubPullRequestActionOpened      = "opened"
	githubPullRequestActionEdited      = "edited"
	githubPullRequestActionSynchronize = "synchronize"
)

func WebhookHandler(prefix string) http.Handler {
	return http.StripPrefix(
		fmt.Sprintf("/%s/", strings.Trim(prefix, "/")),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// Switch action by header
			switch r.Header.Get("X-Github-Event") {
			case githubEventNamePullRequest:
				// Accept PullRequest event
				var evt entity.GithubPullRequestEvent
				if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					io.WriteString(w, "Failed to decode github webhook body to JSON "+err.Error())
					return
				}
				// Get and parse rule from destination repository
				rr, err := github.Content(r.Context(), evt.ContentURL("/.doorkeeper.yml"))
				if err != nil {
					rr = rule.DefaultRule
				}

				// switch actions by action
				switch evt.Action {

				// When new pullrequest has been opened, run validate and factory relates note
				case githubPullRequestActionOpened:
					go validatePullRequest(evt, rr)
					if ok, _ := rr.MatchBranch(evt.BaseBranch()); ok {
						go factoryRelaseNotes(evt)
					}

					// When pullrequest has been edited, only runs validate
				case githubPullRequestActionEdited:
					go validatePullRequest(evt, rr)

					// When pullrequest has been synchronized, only runs factory release notes
				case githubPullRequestActionSynchronize:
					if ok, _ := rr.MatchBranch(evt.BaseBranch()); ok {
						go factoryRelaseNotes(evt)
					}
				}
				successResponse(w)
				return
			case githubEventNamePush:
				// Accept push event
				var evt entity.GithubPushEvent
				if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					io.WriteString(w, "Failed to decode github webhook body to JSON "+err.Error())
					return
				}

				// Get and parse rule from destination repository
				rr, err := github.Content(r.Context(), evt.ContentURL("/.doorkeeper.yml"))
				if err != nil {
					rr = rule.DefaultRule
				}

				switch {
				case strings.HasPrefix(evt.Ref, "refs/tags"):
					if ok, _ := rr.MatchTag(strings.TrimPrefix(evt.Ref, "refs/tags/")); ok {
						go processTagPushEvent(evt, rr)
					}
				}
				successResponse(w)
				return
			case githubEventNamePing:
				// Accept Ping event
				successResponse(w)
				return
			}

			// Forbid other events
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, "We don't support event of '"+r.Header.Get("X-Github-Event")+"'")
		}),
	)
}

func successResponse(w http.ResponseWriter) {
	message := "Accepted"

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(message)))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, message)
}
