package github

import (
	"fmt"
	"log"
	"time"

	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Authorization interface {
	Header(r *http.Request, repository string) error
}

type TokenAuth struct {
	token string
}

func (t *TokenAuth) Header(r *http.Request, repository string) error {
	r.Header.Set("Authorization", "token "+t.token)
	return nil
}

type JWTAuth struct {
	appId   string
	privKey *rsa.PrivateKey
	tokens  AccessTokens
}

func (j *JWTAuth) Header(r *http.Request, repository string) error {
	token, err := j.refreshToken(repository)
	if err != nil {
		return err
	}

	r.Header.Set("Authorization", "token "+token.Token)
	return nil
}

func (j *JWTAuth) refreshToken(repository string) (*AccessToken, error) {
	token := j.tokens.Lookup(repository, func(repository string) *AccessToken {
		jwtToken, err := j.createBearerToken()
		if err != nil {
			log.Println(err)
			return nil
		}
		installationId, err := j.getInstallationId(jwtToken, repository)
		if err != nil {
			log.Println(err)
			return nil
		}
		accessToken, err := j.getAccessToken(jwtToken, installationId)
		if err != nil {
			log.Println(err)
			return nil
		}
		return accessToken
	})
	if token == nil {
		return nil, fmt.Errorf("Failed to get token")
	}
	return token, nil
}

func (k *JWTAuth) getInstallationId(token, repository string) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(
		"https://api.github.com/repos/%s/installation",
		repository,
	), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", githubBasicHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var v struct {
		Id int64 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v.Id, nil
}

func (k *JWTAuth) getAccessToken(token string, installationId int64) (*AccessToken, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(
		"https://api.github.com/app/installations/%d/access_tokens",
		installationId,
	), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", githubBasicHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var v struct {
		Token   string `json:"token"`
		Expires string `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	exp, err := time.Parse(time.RFC3339, v.Expires)
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		Token:  v.Token,
		Expire: exp,
	}, nil
}

func (j *JWTAuth) createBearerToken() (string, error) {
	// update JWT Token
	now := time.Now()
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.StandardClaims{
		IssuedAt:  now.Add(-60 * time.Second).Unix(),
		ExpiresAt: now.Add(10 * time.Minute).Unix(),
		Issuer:    j.appId,
	}
	return token.SignedString(j.privKey)
}
