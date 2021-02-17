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

type Rule struct {
	Title       []RuleItem `yaml:"title"`
	Description []RuleItem `yaml:"description"`
}

func (r *Rule) ValidateTitle(title string) error {
	for i := range r.Title {
		if err := r.validate(r.Title[i], title); err != nil {
			return fmt.Errorf("Invalid PR title: %s", err)
		}
	}
	return nil
}

func (r *Rule) ValidateDescription(desc string) error {
	for i := range r.Description {
		if err := r.validate(r.Description[i], desc); err != nil {
			return fmt.Errorf("Invalid PR description: %s", err)
		}
	}
	return nil
}

func (r *Rule) validate(item RuleItem, dat string) error {
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
	Title: []RuleItem{
		{Kind: "blacklist", Values: []string{"fix", "feature", "implement"}},
	},
	Description: []RuleItem{},
}
