package entity

import (
	"fmt"
	"strings"

	"net/url"
)

type GithubPushEvent struct {
	Ref        string                      `json:"ref"`
	Before     string                      `json:"before"`
	After      string                      `json:"after"`
	Repository GithubPullRequestRepository `json:"repository"`
}

func (e GithubPushEvent) CurrentTag() *Tag {
	return ParseTag(e.Ref, e.Before)
}

func (e GithubPushEvent) ContentURL(path string) string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/contents/%s",
		e.Repository.FullName,
		strings.TrimPrefix(path, "/"),
	)
}

func (e GithubPushEvent) CompareURL(cur, prev *Tag) string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/compare/%s",
		e.Repository.FullName,
		url.PathEscape(prev.Raw+"..."+cur.Raw),
	)
}

func (e GithubPushEvent) TagsURL() string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/git/refs/tags",
		e.Repository.FullName,
	)
}

func (e GithubPushEvent) PullRequestURL(sha string) string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/commits/%s/pulls",
		e.Repository.FullName,
		sha,
	)
}
