package rule

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		input string
		rule  RuleItem
	}{
		{
			input: "feat: support user input",
			rule: RuleItem{
				Kind: RuleKindPrefixed,
				Values: []string{
					"fix:",
					"feat:",
				},
			},
		},
		{
			input: "feat: support user input",
			rule: RuleItem{
				Kind: RuleKindRegex,
				Values: []string{
					".+:.+",
				},
			},
		},
		{
			input: "# Why do you need this change?\n\nUser wants this feature",
			rule: RuleItem{
				Kind: RuleKindContains,
				Values: []string{
					"Why do you need this change?",
				},
			},
		},
		{
			input: "fix",
			rule: RuleItem{
				Kind: RuleKindBlackList,
				Values: []string{
					"feat",
				},
			},
		},
	}

	for i, tt := range tests {
		if err := tt.rule.validate(tt.input); err != nil {
			t.Fatalf("Test[%d] validation error: %s", i, err)
		}
	}
}
