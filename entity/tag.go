package entity

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var TagMatcher = regexp.MustCompile(`v?([0-9]+).([0-9]+).([0-9]+)`)

type Tag struct {
	Raw string
	Sig int64
	Sha string
}

func ParseTag(ref, sha string) *Tag {
	raw := strings.TrimPrefix(ref, "refs/tags/")
	match := TagMatcher.FindStringSubmatch(raw)
	if match == nil {
		return nil
	}
	major, _ := strconv.Atoi(match[1])
	minor, _ := strconv.Atoi(match[2])
	patch, _ := strconv.Atoi(match[3])

	sig, err := strconv.ParseInt(fmt.Sprintf("%d%04d%04d", major, minor, patch), 10, 64)
	if err != nil {
		return nil
	}

	return &Tag{
		Raw: raw,
		Sig: sig,
		Sha: sha,
	}
}

type Tags []*Tag

func (t Tags) PreviousTag(cur *Tag) *Tag {
	sort.SliceStable(t, func(i, j int) bool {
		return t[i].Sig > t[j].Sig
	})

	for i := range t {
		if t[i].Sig > cur.Sig {
			continue
		}
		return t[i]
	}
	return nil
}
