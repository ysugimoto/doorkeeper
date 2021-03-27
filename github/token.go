package github

import (
	"time"
)

type AccessToken struct {
	Token  string
	Expire time.Time
}

type AccessTokens map[string]*AccessToken

func (a AccessTokens) Lookup(repository string, getter func(repository string) *AccessToken) *AccessToken {
	v, ok := a[repository]
	if !ok || v == nil {
		a[repository] = getter(repository)
		return a[repository]
	} else if time.Now().UTC().After(v.Expire) {
		a[repository] = getter(repository)
		return a[repository]
	}
	return v
}
