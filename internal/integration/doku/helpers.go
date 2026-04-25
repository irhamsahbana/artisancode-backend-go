package doku

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *dokuClient) validate() error {
	switch {
	case c.clientID == "":
		return errors.New("doku client id is required")
	case c.secretKey == "":
		return errors.New("doku secret key is required")
	case c.baseURL == "":
		return errors.New("doku base url is required")
	default:
		return nil
	}
}

func (c *dokuClient) generateSignature(payload []byte, timestamp string, requestID string, targetPath string) string {
	digest := sha256.Sum256(payload)
	digestBase64 := base64.StdEncoding.EncodeToString(digest[:])

	component := fmt.Sprintf(
		"Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		c.clientID,
		requestID,
		timestamp,
		targetPath,
		digestBase64,
	)

	return "HMACSHA256=" + c.sign(component)
}

func (c *dokuClient) generateGetSignature(timestamp string, requestID string, targetPath string) string {
	component := fmt.Sprintf(
		"Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s",
		c.clientID,
		requestID,
		timestamp,
		targetPath,
	)

	return "HMACSHA256=" + c.sign(component)
}

func (c *dokuClient) sign(component string) string {
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(component))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func decodeResponse(resp *http.Response, out any) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("doku api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if len(body) == 0 {
		return errors.New("doku api returned empty response body")
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("unmarshal doku response: %w", err)
	}

	return nil
}
