package entity

type GithubCommit struct {
	Sha string `json:"sha"`
}

type GithubReview struct {
	Body  string `json:"body"`
	Event string `json:"event"`
}

type GithubStatus struct {
	Status      string `json:"state"`
	Context     string `json:"context,omitempty"`
	TargetUrl   string `json:"target_url,omitempty"`
	Description string `json:"description,omitempty"`
}
