package rule

import (
	"fmt"
	"regexp"
	"strings"
)

type RuleKind string

const (
	RuleKindPrefixed  RuleKind = "prefixed"
	RuleKindRegex              = "regexp"
	RuleKindContains           = "contains"
	RuleKindBlackList          = "blacklist"
)

type RuleItem struct {
	Kind   RuleKind `yaml:"kind"`
	Values []string `yaml:"values"`
}

type ValidateRule struct {
	Disable     bool       `yaml:"disable"`
	Title       []RuleItem `yaml:"title"`
	Description []RuleItem `yaml:"description"`
	Branches    []string   `yaml:"branches"`
}

type ReleaseNote struct {
	Disable     bool              `yaml:"disable"`
	Branches    []string          `yaml:"branches"`
	Tags        []string          `yaml:"tags"`
	Integration map[string]string `yaml:"integration"`
}

type Rule struct {
	Validation  ValidateRule `yaml:"validation"`
	ReleaseNote ReleaseNote  `yaml:"relasenote"`
}

func (r *Rule) ValidateTitle(title string) error {
	for i := range r.Validation.Title {
		if err := r.Validation.Title[i].validate(title); err != nil {
			return fmt.Errorf("Invalid PR title: %s", err)
		}
	}
	return nil
}

func (r *Rule) ValidateDescription(desc string) error {
	for i := range r.Validation.Description {
		if err := r.Validation.Description[i].validate(desc); err != nil {
			return fmt.Errorf("Invalid PR description: %s", err)
		}
	}
	return nil
}

func (r *Rule) MatchValidateBranch(branch string) (bool, error) {
	return r.matchBranch(branch, r.Validation.Branches)
}

func (r *Rule) MatchReleaseNoteBranch(branch string) (bool, error) {
	return r.matchBranch(branch, r.ReleaseNote.Branches)
}

func (r *Rule) matchBranch(branch string, targets []string) (bool, error) {
	for i := range r.ReleaseNote.Branches {
		if matched, err := regexp.MatchString(r.ReleaseNote.Branches[i], branch); err != nil {
			return false, err
		} else if matched {
			return true, nil
		}
	}
	return false, nil
}

func (r *Rule) MatchTag(tag string) (bool, error) {
	for i := range r.ReleaseNote.Tags {
		if matched, err := regexp.MatchString(r.ReleaseNote.Tags[i], tag); err != nil {
			return false, err
		} else if matched {
			return true, nil
		}
	}
	return false, nil
}

func (r *Rule) Integration() map[string]string {
	return r.ReleaseNote.Integration
}

func (item RuleItem) validate(dat string) error {
	var matched bool

	switch item.Kind {
	// prefixed rule validates as "OR"
	case RuleKindPrefixed:
		for i := range item.Values {
			if matched = strings.HasPrefix(dat, item.Values[i]); matched {
				break
			}
		}
		if !matched {
			return fmt.Errorf("value is not prefixed any of %s", strings.Join(item.Values, ", "))
		}
		return nil
	// contains rule validates as "OR"
	case RuleKindRegex:
		var err error
		for i := range item.Values {
			matched, err = regexp.MatchString("(?is)"+item.Values[i], dat)
			if err != nil {
				return fmt.Errorf("Failed to compile regexp: %s, error: %w", item.Values[i], err)
			}
			if matched {
				break
			}
		}
		if !matched {
			return fmt.Errorf("value is not matched regexp any of %s", strings.Join(item.Values, ", "))
		}
		return nil
	// contains rule validates as "AND"
	case RuleKindContains:
		for i := range item.Values {
			if matched = strings.Contains(dat, item.Values[i]); !matched {
				return fmt.Errorf("value must contain word/section of %s", strings.Join(item.Values, ", "))
			}
		}
		return nil
	// contains rule validates as "AND"
	case RuleKindBlackList:
		for i := range item.Values {
			if item.Values[i] == strings.TrimSpace(dat) {
				return fmt.Errorf("value contains in blacklist of %s", strings.Join(item.Values, ", "))
			}
		}
		return nil
	default:
		return fmt.Errorf("Unexpected validation kind: %s", item.Kind)
	}
}

var DefaultRule = &Rule{
	Validation: ValidateRule{
		Title: []RuleItem{
			{Kind: "blacklist", Values: []string{"fix", "feature", "implement"}},
		},
		Description: []RuleItem{},
	},
	ReleaseNote: ReleaseNote{
		Branches: []string{"deployment/production"},
	},
}
