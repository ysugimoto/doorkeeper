package handler

import (
	"io"
	"regexp"

	"encoding/json"
	"net/http"

	"github.com/ysugimoto/doorkeeper/entity"
)

const (
	githubEventNamePullRequest = "pull_request"
	githubEventNamePing        = "ping"

	githubPullRequestActionOpened      = "opened"
	githubPullRequestActionEdited      = "edited"
	githubPullRequestActionSynchronize = "synchronize"
)

func WebhookHandler() http.Handler {
	return http.StripPrefix("/webhook/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		releaseBranch := r.URL.Path
		if releaseBranch == "" {
			// deployment/production as default
			releaseBranch = "deployment/production"
		}

		branchRegex, err := regexp.Compile(releaseBranch)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Branch regex compilation failed: "+err.Error())
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
			// switch actions by action
			switch evt.Action {

			// When new pullrequest has been opened, run validate and factory relates note
			case githubPullRequestActionOpened:
				go validatePullRequest(evt)
				if branchRegex.MatchString(evt.BaseBranch()) {
					go factoryRelaseNotes(evt)
				}

			// When pullrequest has been edited, only runs validate
			case githubPullRequestActionEdited:
				go validatePullRequest(evt)

			// When pullrequest has been synchronized, only runs factory release notes
			case githubPullRequestActionSynchronize:
				if branchRegex.MatchString(evt.BaseBranch()) {
					go factoryRelaseNotes(evt)
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
	}))
}

func successResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Accepted")
}
