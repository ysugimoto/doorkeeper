package main

import (
	"errors"
	"log"
	"os"

	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/http"

	"github.com/ysugimoto/doorkeeper/github"
	"github.com/ysugimoto/doorkeeper/handler"
)

func main() {
	port := "9000"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}

	privKey, err := parsePrivateKey(os.Getenv("GITHUB_APP_PRIVATE_KEY"))
	if err != nil {
		log.Fatalln(err)
	}

	h := handler.WebhookHandler("/webhook", github.NewClient(
		http.DefaultClient,
		github.WithJWTToken(os.Getenv("GITHUB_APP_ID"), privKey),
	))
	log.Printf("Server starts on :%s", port)
	if err := http.ListenAndServe(":"+port, h); err != nil {
		log.Fatalln(err)
	}
}

func parsePrivateKey(key string) (*rsa.PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(decoded)
	if block == nil {
		return nil, errors.New("PEM block is nil")
	}

	var privKey *rsa.PrivateKey

	switch block.Type {
	case "RSA PRIVATE KEY":
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	case "PRIAVATE KEY":
		ki, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		var ok bool
		if privKey, ok = ki.(*rsa.PrivateKey); !ok {
			return nil, errors.New("Invalid key: not PKC8")
		}
	default:
		return nil, errors.New("Invalid PEM block type")
	}

	privKey.Precompute()
	if err := privKey.Validate(); err != nil {
		return nil, err
	}

	return privKey, nil
}
