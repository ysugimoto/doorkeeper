package entity

import (
	"fmt"
	"strings"
)

type ReleaseNotes struct {
	PullRequestNumber int
	Repository        string
	Notes             []*ReleaseNote
}

func (r ReleaseNotes) SlackMessage() string {
	return ":robot_face: " + r.generateMessage()
}

func (r ReleaseNotes) GitHubMessage() string {
	return ":robot: " + r.generateMessage()
}

func (r ReleaseNotes) generateMessage() string {
	if len(r.Notes) == 0 {
		return fmt.Sprintf(
			"There are no release notes found on release PullRequest <https://github.com/%s/pull/%d|#%d>",
			r.Repository, r.PullRequestNumber, r.PullRequestNumber,
		)
	}

	notes := make([]string, len(r.Notes))
	for i := range r.Notes {
		notes[i] = r.Notes[i].Format(r.Repository)
	}

	return fmt.Sprintf(
		"Release Note collected succesfully on release PullRequest <https://github.com/%s/pull/%d|#%d>\n\n```\n%s\n```\n\nTake this!",
		r.Repository, r.PullRequestNumber, r.PullRequestNumber, strings.Join(notes, "\n"),
	)
}

type ReleaseNote struct {
	PullRequestNumber int
	Note              string
}

func (r *ReleaseNote) Format(repo string) string {
	return fmt.Sprintf(
		"- <https://github.com/%s/pull/%d|#%d> %s",
		repo, r.PullRequestNumber, r.PullRequestNumber, r.Note,
	)
}
