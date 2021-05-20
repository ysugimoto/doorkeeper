package handler

import (
	"testing"
)

func TestFormatReleaseNoteText(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  `Foobar`,
			expect: `Foobar`,
		},
		{
			input:  `<!-- Put release note here -->Foobar`,
			expect: `Foobar`,
		},
		{
			input: `
			<!-- Put release note here -->
			Foobar
			<!-- /Put release note here -->
			`,
			expect: `Foobar`,
		},
	}

	for i, tt := range tests {
		v := formatReleaseNoteText(tt.input)
		if v != tt.expect {
			t.Fatalf(
				"Test failes for TestFormatReleaseNoteText of %d, expect=%s, actual=%s",
				i, tt.expect, v,
			)
		}
	}
}
