package github

import (
	"crypto/rsa"
)

type Option func(gc *Client)

func WithToken(token string) Option {
	return func(gc *Client) {
		gc.authorization = &TokenAuth{
			token: token,
		}
	}
}

func WithJWTToken(appId string, privKey *rsa.PrivateKey) Option {
	return func(gc *Client) {
		gc.authorization = &JWTAuth{
			appId:   appId,
			privKey: privKey,
			tokens:  AccessTokens{},
		}
	}
}
