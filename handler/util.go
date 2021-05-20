package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	htmlCommentRegex          = regexp.MustCompile(`(?s)<!--\s*([^>]+)?>`)
	unnecessaryCharacterRegex = regexp.MustCompile(`(?s)^[\s\t\r\n]+|[\s\t\r\n]+$`)
)

func successResponse(w http.ResponseWriter) {
	message := "Accepted"

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(message)))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, message)
}

// compares webhook request signature with secret
func compareSignature(r *http.Request, secret string) bool {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		fmt.Println("Failed to read request body")
		return false
	}

	// Rewind request body
	defer func() {
		r.Body = ioutil.NopCloser(buf)
	}()

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(buf.Bytes())
	expected := append([]byte("sha256="), []byte(fmt.Sprintf("%x", mac.Sum(nil)))...)
	return hmac.Equal(expected, []byte(r.Header.Get("X-Hub-Signature-256")))
}

// Integration for slack
func sendToSlack(ctx context.Context, webhookURL, message string) error {
	body, err := json.Marshal(map[string]string{
		"text": message,
	})
	if err != nil {
		return fmt.Errorf("Failed to marshal body: %w", err)
	}

	c, timeout := context.WithTimeout(ctx, 5*time.Second)
	defer timeout()

	req, err := http.NewRequestWithContext(c, http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("Failed to make slack request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to get slack response: %w", err)
	}
	resp.Body.Close()
	return nil
}

func formatReleaseNoteText(txt string) string {
	return unnecessaryCharacterRegex.ReplaceAllString(
		htmlCommentRegex.ReplaceAllString(txt, ""),
		"",
	)
}
