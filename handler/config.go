package handler

import (
	"github.com/ysugimoto/doorkeeper/github"
)

type Config struct {
	client    *github.Client
	appSecret string
	token     string
	prefix    string
}

func (c *Config) Client() *github.Client {
	if c.client != nil {
		return c.client
	}

	// If github client is not initialized, use default with supplied token
	return github.DefaultClient(github.WithToken(c.token))
}

type Option func(c *Config)

func WithClient(client *github.Client) Option {
	return func(c *Config) {
		c.client = client
	}
}

func WithAppSecret(s string) Option {
	return func(c *Config) {
		c.appSecret = s
	}
}

func WithToken(t string) Option {
	return func(c *Config) {
		c.token = t
	}
}

func WithPrefix(p string) Option {
	return func(c *Config) {
		c.prefix = p
	}
}
