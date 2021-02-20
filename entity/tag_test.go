package entity

import (
	"testing"
)

func TestTag(t *testing.T) {
	tests := []struct {
		ref    string
		expect int64
	}{
		{
			ref:    "refs/tags/v1.0.0",
			expect: 100000000,
		},
		{
			ref:    "refs/tags/1.0.0",
			expect: 100000000,
		},
		{
			ref:    "refs/tags/v2.12.8",
			expect: 200120008,
		},
	}

	for i, tt := range tests {
		tag := ParseTag(tt.ref, "examplesha")
		if tag == nil {
			t.Fatalf("Test[%d] parsed tag is nil for ref %s", i, tt.ref)
		}
		if tag.Sig != tt.expect {
			t.Fatalf("Test[%d] siganture unmatch. expect=%d got=%d", i, tt.expect, tag.Sig)
		}
	}
}

func TestPreviousTag(t *testing.T) {
	cur := ParseTag("refs/tags/v1.1.1", "examplesha")

	tests := []struct {
		tags   []string
		expect string
	}{
		{
			tags:   []string{"refs/tags/v1.0.0", "refs/tags/v0.0.10"},
			expect: "v1.0.0",
		},
		{
			tags:   []string{"refs/tags/v1.1.0", "refs/tags/v1.1.0"},
			expect: "v1.1.0",
		},
		{
			tags:   []string{"refs/tags/v0.1.0", "refs/tags/v0.0.10"},
			expect: "v0.1.0",
		},
	}

	for i, tt := range tests {
		var tags Tags
		for j := range tt.tags {
			tags = append(tags, ParseTag(tt.tags[j], "examplesha"))
		}
		previous := tags.PreviousTag(cur)
		if previous == nil {
			t.Fatalf("Test[%d] previous tag is nil", i)
		}
		if previous.Raw != tt.expect {
			t.Fatalf("Test[%d] expected tag unmatch. expect=%s got=%s", i, tt.expect, previous.Raw)
		}
	}
}
