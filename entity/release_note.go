package entity

import (
	"fmt"
	"strings"

	"math/rand"
	"time"
)

const (
	slackEmojiPrefix  = ":robot_face:"
	githubEmojiPrefix = ":robot:"
)

type ReleaseNotes struct {
	PullRequestNumber int
	Repository        string
	Notes             []*ReleaseNote
}

func (r ReleaseNotes) SlackMessage() string {
	if len(r.Notes) == 0 {
		return fmt.Sprintf(
			"%s There are no release notes found on release PullRequest <https://github.com/%s/pull/%d|#%d>",
			slackEmojiPrefix, r.Repository, r.PullRequestNumber, r.PullRequestNumber,
		)
	}

	notes := make([]string, len(r.Notes))
	for i := range r.Notes {
		notes[i] = r.Notes[i].SlackFormat(r.Repository)
	}

	return fmt.Sprintf(
		"%s Release Note collected succesfully on release PullRequest <https://github.com/%s/pull/%d|#%d>\n\n```\n%s\n```\n\n\n%s",
		slackEmojiPrefix, r.Repository, r.PullRequestNumber, r.PullRequestNumber, strings.Join(notes, "\n"), randomMessage(),
	)
}

func (r ReleaseNotes) GitHubMessage() string {
	if len(r.Notes) == 0 {
		return fmt.Sprintf(
			"%s There are no release notes found on release PullRequest <https://github.com/%s/pull/%d|#%d>",
			githubEmojiPrefix, r.Repository, r.PullRequestNumber, r.PullRequestNumber,
		)
	}

	notes := make([]string, len(r.Notes))
	for i := range r.Notes {
		notes[i] = r.Notes[i].GitHubFormat(r.Repository)
	}

	return fmt.Sprintf(
		"%s Release Note collected succesfully on release PullRequest #%d\n\n%s\n\n\n%s",
		githubEmojiPrefix, r.PullRequestNumber, strings.Join(notes, "\n"), randomMessage(),
	)
}

type ReleaseNote struct {
	PullRequestNumber int
	Note              string
}

func (r *ReleaseNote) SlackFormat(repo string) string {
	return fmt.Sprintf(
		"- <https://github.com/%s/pull/%d|#%d> %s",
		repo, r.PullRequestNumber, r.PullRequestNumber, r.Note,
	)
}

func (r *ReleaseNote) GitHubFormat(repo string) string {
	return fmt.Sprintf("- #%d %s", r.PullRequestNumber, r.Note)
}

var randoms = []string{
	"Take this!",
	"Hope these releases go well :heart:",
	"Great progress :+1:",
	"Great, team!",
}

func randomMessage() string {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(randoms))

	return randoms[idx]
}
