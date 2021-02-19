package entity

import (
	"fmt"
	"strings"

	"net/url"
)

type GithubPullRequestEvent struct {
	Action      string                      `json:"action"`
	Number      int                         `json:"number"`
	PullRequest GithubPullRequest           `json:"pull_request"`
	Repository  GithubPullRequestRepository `json:"repository"`
}

type GithubPullRequest struct {
	Links      map[string]GithubPullRequestLink `json:"_links"`
	Repository GithubPullRequestRepository      `json:"repo"`
	Number     int                              `json:"number"`
	Title      string                           `json:"title"`
	Body       string                           `json:"body"`
	Head       GithubPullRequestBranch          `json:"head"`
	Base       GithubPullRequestBranch          `json:"base"`
}

type GithubPullRequestLink struct {
	Href string `json:"href"`
}

type GithubPullRequestRepository struct {
	FullName     string                           `json:"full_name"`
	Owner        GithubPullRequestRepositoryOwner `json:"owner"`
	MasterBranch string                           `json:"master_branch"`
}

type GithubPullRequestRepositoryOwner struct {
	Login string `json:"login"`
}

type GithubPullRequestBranch struct {
	Ref string `json:"ref"`
}

func (e GithubPullRequestEvent) HeadBranch() string {
	return e.PullRequest.Head.Ref
}

func (e GithubPullRequestEvent) BaseBranch() string {
	return e.PullRequest.Base.Ref
}

func (e GithubPullRequestEvent) StatusURL() string {
	v, ok := e.PullRequest.Links["statuses"]
	if !ok {
		return ""
	}
	return v.Href
}

func (e GithubPullRequestEvent) ReviewURL() string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/pulls/%d/reviews",
		e.Repository.FullName,
		e.Number,
	)
}

func (e GithubPullRequestEvent) CompareURL() string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/compare/%s",
		e.Repository.FullName,
		url.PathEscape(e.BaseBranch()+"..."+e.HeadBranch()),
	)
}

func (e GithubPullRequestEvent) PullRequestURL(sha string) string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/commits/%s/pulls",
		e.Repository.FullName,
		sha,
	)
}

func (e GithubPullRequestEvent) ContentURL(path string) string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/contents/%s",
		e.Repository.FullName,
		strings.TrimPrefix(path, "/"),
	)
}

func (e GithubPullRequestEvent) RepositoryPath() string {
	return e.Repository.FullName
}
